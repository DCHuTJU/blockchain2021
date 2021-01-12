package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/w3liu/consensus/libs/gobio"
	"github.com/w3liu/consensus/log"
	"github.com/w3liu/consensus/p2p/conn"
	"github.com/w3liu/consensus/types"
	"go.uber.org/zap"
	"net"
	"pbft_consenus/config"
	"pbft_consenus/consensus"
	"sort"
	"sync"
	"time"
)

// pbft 过程中某一个节点的状态
type State struct {
	address        *types.Address
	proposer       *types.Address
	validators     []*types.Address
	mconn          map[string]*conn.MConnection
	blockTicker    *time.Ticker
	lastBlock      Block
	currentBlock   Block
	step           int
	lock           sync.Mutex
	voteCache      map[string]*types.VoteMessage
	preCommitCache map[string]*types.PreCommitMessage
	commitCache    map[string]*types.CommitMessage
	validationPool []types.Result
	voteResultCache map[string]int
}

var (
	Sequence int64 = 0
	Alpha float64 = 0.10
)
func NewState(cfg *config.Config) *State {
	address, err := types.NewAddress(cfg.Peer.Address)
	if err != nil {
		panic(err)
	}

	if len(cfg.Peer.Seeds) == 0 {
		panic("len(cfg.Peer.Seeds) == 0")
	}

	proposer, err := types.NewAddress(cfg.Peer.Seeds[0])

	if err != nil {
		panic(err)
	}

	seeds := make([]*types.Address, 0)

	for _, seed := range cfg.Peer.Seeds {
		addr, err := types.NewAddress(seed)
		if err != nil {
			panic(err)
		}
		if addr.Id == address.Id {
			continue
		}
		seeds = append(seeds, addr)
	}

	return &State {
		proposer:       proposer,
		address:        address,
		validators:     seeds,
		mconn:          make(map[string]*conn.MConnection),
		blockTicker:    time.NewTicker(time.Second * 1),
		voteCache:      make(map[string]*types.VoteMessage),
		preCommitCache: make(map[string]*types.PreCommitMessage),
	}
}

func (s *State) Start() {
	go func() {
		s.dialPeers()
	}()

	go func() {
		s.accept()
	}()

	for {
		select {
		case <-s.blockTicker.C:
			s.Propose()
		}
	}
}

func (s *State) Stop() {

}

func (s *State) accept() {
	ln, err := net.Listen("tcp", s.address.ToIpPortString())
	if err != nil {
		panic(err)
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("ln.Accept() error", err)
			continue
		}

		go func(c net.Conn) {
			server, err := s.wrapConn(c)
			if err != nil {
				fmt.Println("s.wrapConn(c) error", err)
				return
			}
			err = server.OnStart()
			if err != nil {
				fmt.Println("server.OnStart()", err)
				return
			}
		}(c)
	}
}

func (s *State) dialPeers() {
	for _, addr := range s.validators {
		if _, ok := s.mconn[addr.Id]; !ok {
			// 尝试拨号
			c, err := s.dial(addr)
			if err != nil {
				log.Info("s.dial(addr) error", zap.Error(err))
				continue
			}

			err = c.OnStart()
			if err != nil {
				log.Info("c.OnStart()", zap.Error(err))
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func (s *State) dial(addr *types.Address) (*conn.MConnection, error) {
	c, err := net.DialTimeout("udp", addr.ToIpPortString(), time.Second * 3)
	if err != nil {
		return nil, err
	}

	return s.wrapConn(c)
}

func (s *State) wrapConn(c net.Conn) (*conn.MConnection, error) {
	peer, err := s.handshake(c, time.Second * 3, &types.NodeInfoMessage{
		Id: s.address.Id,
		Ip: s.address.Ip.String(),
		Port: s.address.Port,
	})

	if err != nil {
		log.Error("handshake error", zap.Error(err))
		_ = c.Close()
		return nil, err
	}

	if _, ok := s.mconn[peer.Id]; ok {
		log.Warn("connection is existed", zap.Any("peer", peer))
		return nil, err
	}

	onError := func(r interface{}) {
		// 移除节点
		if _, ok := s.mconn[peer.Id]; ok {
			delete(s.mconn, peer.Id)
		}
		log.Error("onError", zap.Any("r", r))
	}

	chDescs := []*conn.ChannelDescriptor{{ID: 0x01, Priority: 1, SendQueueCapacity: 1}}
	server := conn.NewMConnection(c, chDescs, s.ReceiveMsg, onError)
	s.mconn[peer.Id] = server
	return server, nil
}

func (s *State) handshake(c net.Conn, timeout time.Duration, nodeInfo *types.NodeInfoMessage) (*types.NodeInfoMessage, error) {
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	var (
		errc         = make(chan error, 2)
		peerNodeInfo = &types.NodeInfoMessage{}
		ourNodeInfo  = nodeInfo
	)
	go func(errc chan<- error, c net.Conn) {
		_, err := gobio.NewWriter(c).WriteMsg(ourNodeInfo)
		errc <- err
	}(errc, c)

	go func(errc chan<- error, c net.Conn) {
		err := gobio.NewReader(c).ReadMsg(peerNodeInfo)
		errc <- err
	}(errc, c)

	for i := 0; i < cap(errc); i++ {
		err := <-errc
		if err != nil {
			return nil, err
		}
	}
	return peerNodeInfo, c.SetDeadline(time.Time{})
}

func (s *State) isProposer() bool {
	return s.proposer.Id == s.address.Id
}

func (s *State) sortRequest() types.Result {
	sort.Slice(s.validationPool, func(i, j int) bool {
		return s.validationPool[i].Reward > s.validationPool[i].Reward
	})
	return s.validationPool[0]
}

// 对所有结果进行排序
func (s *State) checkRequest(results []types.Result) types.Result {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Reward > results[j].Reward
	})
	return results[0]
}

func (s *State) checkBlock() bool {
	return true
}

func (s *State) checkMajor23(cnt int) bool {
	n := len(s.validators) + 1
	return (cnt * 1000 / n) > 2*1000/3
}

func (s *State) ReceiveMsg(chID byte, msgBytes []byte) {
	//log.Info("ReceiveMsg:", zap.Any("chID", chID), zap.String("msg", string(msgBytes)))
	msgInfo := &types.MessageInfo{}
	err := json.Unmarshal(msgBytes, msgInfo)
	if err != nil {
		log.Error("json.Unmarshal", zap.Error(err))
		return
	}
	switch msgInfo.MsgType {
	case types.Propose:
		proposeMsg := &types.ProposeMessage{}
		err := json.Unmarshal([]byte(msgInfo.MsgContent), proposeMsg)
		if err != nil {
			log.Error("json.Unmarshal", zap.Error(err))
			return
		}
		s.ReceivePropose(proposeMsg)
	case types.Vote:
		voteMsg := &types.VoteMessage{}
		err := json.Unmarshal([]byte(msgInfo.MsgContent), voteMsg)
		if err != nil {
			log.Error("json.Unmarshal", zap.Error(err))
			return
		}
		s.ReceiveVote(voteMsg)
	case types.PreCommit:
		preCommitMsg := &types.PreCommitMessage{}
		err := json.Unmarshal([]byte(msgInfo.MsgContent), preCommitMsg)
		if err != nil {
			log.Error("json.Unmarshal", zap.Error(err))
			return
		}
		s.ReceivePreCommit(preCommitMsg)

	}
}



// 请求发布阶段
func (s *State) ProposeRequest(request types.RequestSet) {
	var tmp types.RequestSet
	// 共识开始，分发所有请求
	if s.step == types.Pending && s.isProposer() {
		tmp = request
	} else {
		if s.isProposer() {
			log.Warn("Propose", zap.Any("step", s.step))
		}
		return
	}

	// 检查节点数量是否大于 2/3
	if !s.checkMajor23(len(s.mconn) + 1) {
		log.Warn("connected peer number is lower than 2/3")
		return
	}
	// 全局分配唯一序号
	Sequence += 1
	msg := &types.RequestMessage{
		RequestSet: tmp,
		TimeStamp: time.Now(),
		Number: Sequence,
	}
	// 向所有其他节点发送当前信息
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	// 状态发生变化，变成投票阶段
	s.step = types.VoteRequest
}

// 请求接受阶段
func (s *State) ReceiveProposeRequest(o *types.RequestMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 在接受到请求后，对请求进行计算，获得结果，并广播出去
	if s.step == types.Pending {
		rlt := consensus.Consensus()
		// 将结果保存在自己的 验证池 中
		s.validationPool = append(s.validationPool, rlt)
		s.BroadcastResult(rlt)
	}
}

func (s *State) BroadcastResult(rlt types.Result) {
	msg := &types.BroadcastResultMessage{
		Result: rlt,
	}
	// 将结果广播出去
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	s.step = types.BroadcastResult
}

func (s *State) ReceiveBroadcastResult(o *types.BroadcastResultMessage) {

	if s.step == types.BroadcastResult {
		duration := 10 * time.Second
		timeout := Alpha * duration.Seconds()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
		defer cancel()

		done := make(chan struct{}, 1)

		// 在 1 + alpha 的时间内接收其他节点发来的信息
		go func() {
			s.ReceiveBroadcastResult(o)
			s.validationPool = append(s.validationPool, o.Result)
			done <- struct{}{}
		}()

		select {
		case <-done:
			fmt.Println("done")
			goto endfor
		case <-ctx.Done():
			fmt.Println("main", ctx.Err())
		}

		//time.Sleep(time.Duration(timeout))
		//go func(o *types.BroadcastResultMessage) {
		//
		//}(o)
	}
	endfor:
		rlt := s.sortRequest()
		s.VoteRequest(rlt)
}

func (s *State) VoteRequest(result types.Result) {
	msg := &types.VoteReqeustMessage{
		Result: result.Reward,
		PeerInfo: result.Id,
	}
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	s.step = types.VoteRequest
}

func (s *State) VoteRequestReceive(o *types.VoteReqeustMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.step <= types.Vote {
		// 若当前投票信息不存在的话
		if _, ok := s.voteResultCache[o.PeerInfo]; !ok {
			// 就将信息添加到其中
			s.voteResultCache[o.PeerInfo]++
			// 如果满足大多数的话就成为 proposer
			if s.checkMajor23(s.voteResultCache[s.address.Id]) {
				s.proposer = s.address
				s.Propose()
			}
		}
	}
}

// 区块更新提交过程
func (s *State) Propose() {
	if s.step == types.Pending && s.isProposer() {
		currHeight := s.lastBlock.Height + 1
		s.currentBlock = Block{
			Height: currHeight,
			Data:   fmt.Sprintf("This is a block data, height is %d.", currHeight),
		}
	} else {
		if s.isProposer() {
			log.Warn("Propose", zap.Any("step", s.step))
		}
		return
	}

	if !s.checkMajor23(len(s.mconn) + 1) {
		log.Warn("connected peer number is lower than 2/3")
		return
	}
	msg := &types.ProposeMessage{
		Height:    s.currentBlock.Height,
		Validator: s.address.Id,
		Data:      s.currentBlock.Data,
		Signer:    s.address.Id,
	}
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	s.step = types.Vote
}

func (s *State) Vote() {
	// 封装信息
	msg := &types.VoteMessage{
		Height:    s.currentBlock.Height,
		Validator: s.address.Id,
	}
	// 广播
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	// 状态转换为 投票 状态
	s.step = types.Vote
}

func (s *State) PreCommit() {
	msg := &types.PreCommitMessage{
		Height:    s.currentBlock.Height,
		Validator: s.address.Id,
	}
	for _, c := range s.mconn {
		if b := c.Send(0x1, types.NewMessageInfo(msg).GetData()); !b {
			log.Error("send msg failed")
			return
		}
	}
	s.step = types.PreCommit
}

func (s *State) Commit() {
	log.Info("block commit", zap.Any("height", s.currentBlock.Height), zap.Any("data", s.currentBlock.Data))
	s.lastBlock = s.currentBlock
	s.currentBlock = Block{}
	for k, _ := range s.voteCache {
		delete(s.voteCache, k)
	}
	for k, _ := range s.preCommitCache {
		delete(s.preCommitCache, k)
	}
	s.step = types.Pending
}

// 接受 propose 请求
func (s *State) ReceivePropose(o *types.ProposeMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.step == types.Pending {
		// 当前区块即为接收到的区块
		s.currentBlock = Block{
			Height: o.Height,
			Data:   o.Data,
		}
		//rlt := s.checkBlock(results)
		//s.Vote(rlt)
		// 如果区块合理，则发送投票
		if s.checkBlock() {
			s.Vote()
		} else {
			// 否则仍然保持原有区块
			s.currentBlock = s.lastBlock
		}
	}
}

func (s *State) ReceiveVote(o *types.VoteMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.step <= types.Vote {
		// 若当前投票信息不存在的话
		if _, ok := s.voteCache[o.Validator]; !ok {
			// 就将信息添加到其中
			s.voteCache[o.Validator] = o
		}
		var voteCnt int
		if s.isProposer() {
			voteCnt = len(s.voteCache) + 1
		} else {
			voteCnt = len(s.voteCache) + 2
		}
		// 如果满足大多数的话，就预提交
		if s.checkMajor23(voteCnt) {
			s.PreCommit()
		}
	}
}

func (s *State) ReceivePreCommit(o *types.PreCommitMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.step == types.PreCommit {
		if _, ok := s.preCommitCache[o.Validator]; !ok {
			s.preCommitCache[o.Validator] = o
		}
		var pCnt = len(s.preCommitCache) + 1
		if s.checkMajor23(pCnt) {
			s.Commit()
		}
	}
}