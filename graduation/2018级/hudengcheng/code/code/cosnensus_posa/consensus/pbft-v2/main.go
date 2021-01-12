package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	mathRand "math/rand"
)

const nodeCount = 4

//客户端的监听地址
var clientAddr = "127.0.0.1:8888"

// 用于统计所有时长
var timePrePrepare, timePrepare, timeCommit, timeReply, timeTotal int64

//节点池，主要用来存储监听地址
var nodeTable map[string]string

// 上一轮共识的领导者，初始化默认为N0
var lastLeader = nodeTable["N0"]


func main() {
	//为四个节点生成公私钥
	genRsaKeys()

	nodeTable = ReadFromCSV()
	if len(os.Args) != 2 {
		log.Panic("输入的参数有误！")
	}
	nodeID := os.Args[1]
	if nodeID == "client" {
		clientSendMessageAndListen() //启动客户端程序
	} else if addr, ok := nodeTable[nodeID]; ok {
		p := NewPBFT(nodeID, addr)
		go p.tcpListen() //启动节点
	} else {
		log.Fatal("无此节点编号！")
	}
	select {}
}

// 客户端 client
func clientSendMessageAndListen() {
	//开启客户端的本地监听（主要用来接收节点的reply信息）
	go clientTcpListen()
	fmt.Printf("客户端开启监听，地址：%s\n", clientAddr)

	fmt.Println(" ---------------------------------------------------------------------------------")
	fmt.Println("|  已进入PBFT测试Demo客户端，请启动全部节点后再发送消息！ :)  |")
	fmt.Println(" ---------------------------------------------------------------------------------")
	fmt.Println("请在下方输入要存入节点的信息：")
	// 首先通过命令行获取用户输入
	stdReader := bufio.NewReader(os.Stdin)
	for {
		data, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		r := new(Request)
		r.Timestamp = time.Now().UnixNano()
		r.ClientAddr = clientAddr
		r.Message.ID = getRandom()
		// 消息内容就是用户的输入
		// 输入内容就是 RequestStorage
		r.Message.Content = strings.TrimSpace(data)
		br, err := json.Marshal(r)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(br))
		content := jointMessage(cPreRequest, br)
		// 客户端的信息发送至上一轮的共识领导者
		// 此处应该设计一个方法，传统的是轮询的方法作为主节点
		tcpDial(content, nodeTable["N0"])
	}
}

//返回一个十位数的随机数，作为msgid
func getRandom() int {
	x := big.NewInt(10000000000)
	for {
		result, err := rand.Int(rand.Reader, x)
		if err != nil {
			log.Panic(err)
		}
		if result.Int64() > 1000000000 {
			return int(result.Int64())
		}
	}
}

// cmd
// 请求部分
type StorageRequest struct {
	// 文件大小
	Size int
	// MD5 校验码
	MD5  []byte
}

type StorageReqeustMessage struct {
	Timestamp int64
	Requests  []StorageRequest
}

type PreRequest struct {
	// 分配一个序号
	N int
	// 请求信息摘要
	Digest string
	// 请求信息实体
	Content StorageReqeustMessage
	Sign  []byte
}

type RequestVoteMessage struct {
	// 反馈的是哪一个请求信息
	N int
	Digest string
	Result string
	Entropy float64
	NodeId string
	Sign   []byte
	Timestamp int64
}

type RequestCommitMessage struct {
	N int
	Digest string
	// 支持哪个节点
	NodeId string
	Sign []byte
}

type BecomeLeader struct {
	Leader string
	Sign   []byte
}

// 共识部分
// <REQUEST,o,t,c>
type Request struct {
	Message
	Timestamp int64
	//相当于clientID
	ClientAddr string
}

//<<PRE-PREPARE,v,n,d>,m>
type PrePrepare struct {
	RequestMessage Request
	Digest         string
	SequenceID     int
	Sign           []byte
}

//<PREPARE,v,n,d,i>
type Prepare struct {
	Digest     string
	SequenceID int
	NodeID     string
	Sign       []byte
}

//<COMMIT,v,n,D(m),i>
type Commit struct {
	Digest     string
	SequenceID int
	NodeID     string
	Sign       []byte
}

//<REPLY,v,t,c,i,r>
type Reply struct {
	MessageID int
	NodeID    string
	Result    bool
}

type Message struct {
	Content string
	ID      int
}

const prefixCMDLength = 15

type command string

const (
	cPreRequest command = "prerequest"
	cReqeustVote command = "requestvote"
	cRequestCommit command = "requestcommit"
	cBecomeLeader command = "becomeleader"
	cRequest    command = "request"
	cPrePrepare command = "preprepare"
	cPrepare    command = "prepare"
	cCommit     command = "commit"
)

// 默认前十五位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

// 默认前十五位为命令名称
func splitMessage(message []byte) (cmd string, content []byte) {
	cmdBytes := message[:prefixCMDLength]
	newCMDBytes := make([]byte, 0)
	for _, v := range cmdBytes {
		if v != byte(0) {
			newCMDBytes = append(newCMDBytes, v)
		}
	}
	cmd = string(newCMDBytes)
	content = message[prefixCMDLength:]
	return
}

//对消息详情进行摘要
func getDigest(request Request) string {
	b, err := json.Marshal(request)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	// 进行十六进制字符串编码
	// 长度为 64 位
	return hex.EncodeToString(hash[:])
}

// 对信息详情进行摘要
func (m StorageReqeustMessage) getDigest() string {
	b, err := json.Marshal(m)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	// 进行十六进制字符串编码
	// 长度为 64 位
	return hex.EncodeToString(hash[:])
}


// 核心 pbft
//本地消息池（模拟持久化层），只有确认提交成功后才会存入此池
var localMessagePool = []Message{}


type node struct {
	//节点ID
	nodeID string
	//节点监听地址
	addr string
	//RSA私钥
	rsaPrivKey []byte
	//RSA公钥
	rsaPubKey []byte
}

type pbft struct {
	//节点信息
	node node
	//每笔请求自增序号
	sequenceID int
	//锁
	lock sync.Mutex
	//
	requestPool map[string]StorageReqeustMessage
	//临时消息池，消息摘要对应消息本体
	messagePool map[string]Request
	//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	prePareConfirmCount map[string]map[string]bool
	//存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	commitConfirmCount map[string]map[string]bool
	//该笔消息是否已进行Commit广播
	isCommitBordcast map[string]bool
	//该笔消息是否已对客户端进行Reply
	isReply map[string]bool
	// 本地投票信息接收池，在有限时间内接收到的信息会被存入此池
	voteMessagePool []RequestVoteMessage
}

func NewPBFT(nodeID, addr string) *pbft {
	p := new(pbft)
	p.node.nodeID = nodeID
	p.node.addr = addr
	p.node.rsaPrivKey = p.getPivKey(nodeID) //从生成的私钥文件处读取
	p.node.rsaPubKey = p.getPubKey(nodeID)  //从生成的私钥文件处读取
	p.sequenceID = 0
	p.requestPool = make(map[string]StorageReqeustMessage)
	p.messagePool = make(map[string]Request)
	p.prePareConfirmCount = make(map[string]map[string]bool)
	p.commitConfirmCount = make(map[string]map[string]bool)
	p.isCommitBordcast = make(map[string]bool)
	p.isReply = make(map[string]bool)
	p.voteMessagePool = make([]RequestVoteMessage, 0)
	return p
}

func (p *pbft) handleRequest(data []byte) {
	//切割消息，根据消息命令调用不同的功能
	cmd, content := splitMessage(data)
	switch command(cmd) {
	case cPreRequest:
		p.handleStorageRequest(content)
	case cReqeustVote:
		p.handleRequestAndVote(content)
	case cRequestCommit:
		p.handleRequestCommit(content)
	case cBecomeLeader:
		p.handleBecomeLeader(content)
	//case cRequest:
	//	p.handleClientRequest(content)
	case cPrePrepare:
		p.handlePrePrepare(content)
	case cPrepare:
		p.handlePrepare(content)
	case cCommit:
		p.handleCommit(content)
	}
}

var newLeader = false
var curLeader = ""

// posa
// 由上一轮的领导节点处理客户端发来的请求
func (p *pbft) handleStorageRequest(content []byte) {
	fmt.Println("主节点已接收到客户端发来的request ...")
	start := time.Now().UnixNano()
	// time.Sleep(RandomDelayGenerator())
	//使用json解析出Request结构体
	r := new(StorageReqeustMessage)
	err := json.Unmarshal(content, r)
	if err != nil {
		log.Panic(err)
	}
	//添加信息序号
	p.sequenceIDAdd()
	//获取消息摘要
	digest := r.getDigest()
	fmt.Println("已将request存入临时消息池")
	//存入临时消息池
	p.requestPool[digest] = *r
	//主节点对消息摘要进行签名
	digestByte, _ := hex.DecodeString(digest)
	signInfo := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
	//拼接成PrePrepare，准备发往follower节点
	pp := PreRequest{
		Digest: digest,
		Sign: signInfo,
		N: p.sequenceID,
		Content: *r,
	}
	b, err := json.Marshal(pp)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("正在向其他节点进行进行PreRequest广播 ...")
	// 进行PrePrepare广播
	p.broadcast(cPreRequest, b)
	fmt.Println("PreRequest广播完成")
	end := time.Now().UnixNano()
	timePrePrepare = end - start
	fmt.Println("PreRequest时长:", end - start)
}

// 根据接收到的请求进行运算，获得结果
func (p *pbft) handleRequestAndVote(content []byte) {
	fmt.Println("本节点已接收到主节点发来的PreRequest...")
	start := time.Now().UnixNano()
	pp := new(PreRequest)
	err := json.Unmarshal(content, pp)
	if err != nil {
		log.Panic(err)
	}
	// 处理 content 里面的所有请求，并进行相应的计算
	tmp := pp.Content
	entropy, result := DealWithContent(tmp)

	// 获取节点的公钥，用于数字签名验证
	primaryNodePubKey := p.getPubKey(lastLeader)
	digestByte, _ := hex.DecodeString(pp.Digest)
	if digest := pp.Content.getDigest(); digest != pp.Digest {
		fmt.Println("信息摘要对不上，拒绝进行prepare广播")
	} else if p.sequenceID+1 != pp.N {
		fmt.Println("消息序号对不上，拒绝进行prepare广播")
	} else if !p.RsaVerySignWithSha256(digestByte, pp.Sign, primaryNodePubKey) {
		fmt.Println("主节点签名验证失败！,拒绝进行prepare广播")
	} else {
		//序号赋值
		p.sequenceID = pp.N
		//将信息存入临时消息池
		fmt.Println("已将消息存入临时节点池")
		// 节点使用私钥对其签名
		sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
		// 拼接成VoteMessage
		pre := RequestVoteMessage{
			Digest: pp.Digest,
			Sign: sign,
			N: pp.N,
			Entropy: entropy,
			Result: result,
			NodeId: p.node.nodeID,
			Timestamp: time.Now().UnixNano(),
		}
		// pre := Prepare{pp.Digest, pp.SequenceID, p.node.nodeID, sign}
		p.voteMessagePool = append(p.voteMessagePool, pre)
		bPre, err := json.Marshal(pre)
		if err != nil {
			log.Panic(err)
		}
		//进行准备阶段的广播
		fmt.Println("正在进行RequestVote广播 ...")
		p.broadcast(cReqeustVote, bPre)
		fmt.Println("RequestVote广播完成")
	}
	end := time.Now().UnixNano()
	timePrepare = end - start
	fmt.Println("RequestVote时长:", timePrepare)
}

func DealWithContent(content StorageReqeustMessage) (entropy float64, result string) {
	return 0, ""
}

// 同时需要接收多个节点发来的信息
func (p *pbft) handleRequestCommit(content []byte) {
	start := time.Now().UnixNano()
	T := time.Now().UnixNano()
	//使用json解析出Prepare结构体
	pre := new(RequestVoteMessage)
	err := json.Unmarshal(content, pre)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的Vote... \n", pre.NodeId)
	//获取消息源节点的公钥，用于数字签名验证
	MessageNodePubKey := p.getPubKey(pre.NodeId)
	digestByte, _ := hex.DecodeString(pre.Digest)
	if _, ok := p.messagePool[pre.Digest]; !ok {
		fmt.Println("当前临时消息池无此摘要，拒绝执行commit广播")
	} else if p.sequenceID != pre.N {
		fmt.Println("消息序号对不上，拒绝执行commit广播")
	} else if !p.RsaVerySignWithSha256(digestByte, pre.Sign, MessageNodePubKey) {
		fmt.Println("节点签名验证失败！,拒绝执行commit广播")
	} else {
		p.setPrePareConfirmMap(pre.Digest, pre.NodeId, true)
		count := 0
		for range p.prePareConfirmCount[pre.Digest] {
			count++
			if count == 1 {
				T = pre.Timestamp
			}
		}
		p.voteMessagePool = append(p.voteMessagePool, *pre)

		// 如果在 1 + alpha 时间内接收到信息
		p.lock.Lock()
		//获取消息源节点的公钥，用于数字签名验证
		if time.Now().UnixNano() - T >= 2000000000 {
			tmp := p.voteMessagePool[:count]
			fmt.Printf("本节点已收到 %d 个节点(包括本地节点)发来的VoteMessage信息 ...", count)
			// 对tmp中的信息进行排序
			entropySet := make([]RequestVoteMessage, 0)
			for i:=0; i<len(tmp); i++ {
				entropySet = append(entropySet, tmp[i])
			}
			sort.Slice(entropySet, func(i, j int) bool {
				if entropySet[i].Entropy > entropySet[j].Entropy {
					return true
				}
				return false
			})
			//节点使用私钥对其签名
			sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
			/*
			type RequestCommitMessage struct {
				N int
				Digest string
				// 支持哪个节点
				NodeId string
				Sign []byte
			}
			 */
			c := RequestCommitMessage{
				N: pre.N,
				Digest: string(digestByte),
				NodeId: entropySet[0].NodeId,
				Sign: sign,
			}
			bc, err := json.Marshal(c)
			if err != nil {
				log.Panic(err)
			}
			//进行提交信息的广播
			fmt.Println("正在进行RequestCommit广播")
			p.broadcast(cRequestCommit, bc)
			p.isCommitBordcast[pre.Digest] = true
			fmt.Println("RequestCommit广播完成")
		}
		p.lock.Unlock()
	}
	end := time.Now().UnixNano()
	timeCommit = end - start
	fmt.Println("RequestCommit时长:", end - start)
}

// 同时需要接收多个节点发来的投票信息
func (p *pbft) handleBecomeLeader(content []byte) {
	start := time.Now().UnixNano()
	// time.Sleep(RandomDelayGenerator())
	//使用json解析出Prepare结构体
	pre := new(RequestCommitMessage)
	err := json.Unmarshal(content, pre)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的RequestCommit... \n", pre.NodeId)
	//获取消息源节点的公钥，用于数字签名验证
	MessageNodePubKey := p.getPubKey(pre.NodeId)
	digestByte, _ := hex.DecodeString(pre.Digest)
	if _, ok := p.messagePool[pre.Digest]; !ok {
		fmt.Println("当前临时消息池无此摘要，拒绝执行commit广播")
	} else if p.sequenceID != pre.N {
		fmt.Println("消息序号对不上，拒绝执行commit广播")
	} else if !p.RsaVerySignWithSha256(digestByte, pre.Sign, MessageNodePubKey) {
		fmt.Println("节点签名验证失败！,拒绝执行commit广播")
	} else {
		p.setPrePareConfirmMap(pre.Digest, pre.NodeId, true)
		count := 0
		voteCount := 0
		for range p.prePareConfirmCount[pre.Digest] {
			count++
			if pre.NodeId == p.node.nodeID {
				voteCount++
			}
		}

		//如果节点至少收到了2f个prepare的消息（包括自己）,并且没有进行过commit广播，则进行commit广播
		p.lock.Lock()
		//获取消息源节点的公钥，用于数字签名验证
		if voteCount >= int(math.Ceil(0.5 * nodeCount)) && newLeader == false {
			fmt.Println("本节点已收到至少一般以上的节点(包括本地节点)发来的Vote信息 ...")
			//节点使用私钥对其签名
			sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
			c := PrePrepare{Request{

			}, string(digestByte), p.sequenceID, sign}
			bc, err := json.Marshal(c)
			if err != nil {
				log.Panic(err)
			}
			//进行提交信息的广播
			curLeader = p.node.nodeID
			fmt.Println("正在进行request广播")
			p.broadcast(cRequest, bc)
			p.isCommitBordcast[pre.Digest] = true
			fmt.Println("request广播完成")
		}
		p.lock.Unlock()
	}
	end := time.Now().UnixNano()
	timeCommit = end - start
	fmt.Println("Reply时长:", end - start)
}

//处理预准备消息
func (p *pbft) handlePrePrepare(content []byte) {
	fmt.Println("本节点已接收到主节点发来的PrePrepare ...")
	start := time.Now().UnixNano()
	// time.Sleep(RandomDelayGenerator())
	//	//使用json解析出PrePrepare结构体
	pp := new(PrePrepare)
	err := json.Unmarshal(content, pp)
	if err != nil {
		log.Panic(err)
	}
	//获取主节点的公钥，用于数字签名验证
	primaryNodePubKey := p.getPubKey("N0")
	digestByte, _ := hex.DecodeString(pp.Digest)
	if digest := getDigest(pp.RequestMessage); digest != pp.Digest {
		fmt.Println("信息摘要对不上，拒绝进行prepare广播")
	} else if p.sequenceID+1 != pp.SequenceID {
		fmt.Println("消息序号对不上，拒绝进行prepare广播")
	} else if !p.RsaVerySignWithSha256(digestByte, pp.Sign, primaryNodePubKey) {
		fmt.Println("主节点签名验证失败！,拒绝进行prepare广播")
	} else {
		//序号赋值
		p.sequenceID = pp.SequenceID
		//将信息存入临时消息池
		fmt.Println("已将消息存入临时节点池")
		p.messagePool[pp.Digest] = pp.RequestMessage
		//节点使用私钥对其签名
		sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
		//拼接成Prepare
		pre := Prepare{pp.Digest, pp.SequenceID, p.node.nodeID, sign}
		bPre, err := json.Marshal(pre)
		if err != nil {
			log.Panic(err)
		}
		//进行准备阶段的广播
		fmt.Println("正在进行Prepare广播 ...")
		p.broadcast(cPrepare, bPre)
		fmt.Println("Prepare广播完成")
	}
	end := time.Now().UnixNano()
	timePrepare = end - start
	fmt.Println("Prepare时长:", timePrepare)
}

//处理准备消息
func (p *pbft) handlePrepare(content []byte) {
	start := time.Now().UnixNano()
	// time.Sleep(RandomDelayGenerator())
	//使用json解析出Prepare结构体
	pre := new(Prepare)
	err := json.Unmarshal(content, pre)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的Prepare ... \n", pre.NodeID)
	//获取消息源节点的公钥，用于数字签名验证
	MessageNodePubKey := p.getPubKey(pre.NodeID)
	digestByte, _ := hex.DecodeString(pre.Digest)
	if _, ok := p.messagePool[pre.Digest]; !ok {
		fmt.Println("当前临时消息池无此摘要，拒绝执行commit广播")
	} else if p.sequenceID != pre.SequenceID {
		fmt.Println("消息序号对不上，拒绝执行commit广播")
	} else if !p.RsaVerySignWithSha256(digestByte, pre.Sign, MessageNodePubKey) {
		fmt.Println("节点签名验证失败！,拒绝执行commit广播")
	} else {
		p.setPrePareConfirmMap(pre.Digest, pre.NodeID, true)
		count := 0
		for range p.prePareConfirmCount[pre.Digest] {
			count++
		}
		//因为主节点不会发送Prepare，所以不包含自己
		specifiedCount := 0
		if p.node.nodeID == "N0" {
			specifiedCount = nodeCount / 3 * 2
		} else {
			specifiedCount = (nodeCount / 3 * 2) - 1
		}
		//如果节点至少收到了2f个prepare的消息（包括自己）,并且没有进行过commit广播，则进行commit广播
		p.lock.Lock()
		//获取消息源节点的公钥，用于数字签名验证
		if count >= specifiedCount && !p.isCommitBordcast[pre.Digest] {
			fmt.Println("本节点已收到至少2f个节点(包括本地节点)发来的Prepare信息 ...")
			//节点使用私钥对其签名
			sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
			c := Commit{pre.Digest, pre.SequenceID, p.node.nodeID, sign}
			bc, err := json.Marshal(c)
			if err != nil {
				log.Panic(err)
			}
			//进行提交信息的广播
			fmt.Println("正在进行commit广播")
			p.broadcast(cCommit, bc)
			p.isCommitBordcast[pre.Digest] = true
			fmt.Println("commit广播完成")
		}
		p.lock.Unlock()
	}
	end := time.Now().UnixNano()
	timeCommit = end - start
	fmt.Println("Reply时长:", end - start)
}

// 处理提交确认消息
func (p *pbft) handleCommit(content []byte) {
	start := time.Now().UnixNano()
	// time.Sleep(RandomDelayGenerator())
	//使用json解析出Commit结构体
	c := new(Commit)
	err := json.Unmarshal(content, c)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的Commit ... \n", c.NodeID)
	//获取消息源节点的公钥，用于数字签名验证
	MessageNodePubKey := p.getPubKey(c.NodeID)
	digestByte, _ := hex.DecodeString(c.Digest)
	if _, ok := p.prePareConfirmCount[c.Digest]; !ok {
		fmt.Println("当前prepare池无此摘要，拒绝将信息持久化到本地消息池")
	} else if p.sequenceID != c.SequenceID {
		fmt.Println("消息序号对不上，拒绝将信息持久化到本地消息池")
	} else if !p.RsaVerySignWithSha256(digestByte, c.Sign, MessageNodePubKey) {
		fmt.Println("节点签名验证失败！,拒绝将信息持久化到本地消息池")
	} else {
		p.setCommitConfirmMap(c.Digest, c.NodeID, true)
		count := 0
		for range p.commitConfirmCount[c.Digest] {
			count++
		}
		//如果节点至少收到了2f+1个commit消息（包括自己）,并且节点没有回复过,并且已进行过commit广播，则提交信息至本地消息池，并reply成功标志至客户端！
		p.lock.Lock()
		if count >= nodeCount/3*2 && !p.isReply[c.Digest] && p.isCommitBordcast[c.Digest] {
			fmt.Println("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
			//将消息信息，提交到本地消息池中！
			localMessagePool = append(localMessagePool, p.messagePool[c.Digest].Message)
			info := p.node.nodeID + "节点已将msgid:" + strconv.Itoa(p.messagePool[c.Digest].ID) + "存入本地消息池中,消息内容为：" + p.messagePool[c.Digest].Content
			fmt.Println(info)
			fmt.Println("正在reply客户端 ...")
			tcpDial([]byte(info), p.messagePool[c.Digest].ClientAddr)
			p.isReply[c.Digest] = true
			fmt.Println("reply完毕")
		}
		p.lock.Unlock()
	}
	end := time.Now().UnixNano()
	timeReply = end - start
	fmt.Println("Reply时长:", end - start)
	timeTotal = timePrePrepare + timePrepare + timeReply + timeCommit
	fmt.Println("本轮共识共花费时长为: ", timeTotal)
}

// 序号累加
func (p *pbft) sequenceIDAdd() {
	p.lock.Lock()
	p.sequenceID++
	p.lock.Unlock()
}

// 向除自己外的其他节点进行广播
func (p *pbft) broadcast(cmd command, content []byte) {
	for i := range nodeTable {
		if i == p.node.nodeID {
			continue
		}
		message := jointMessage(cmd, content)
		go tcpDial(message, nodeTable[i])
	}
}

//为多重映射开辟赋值
func (p *pbft) setPrePareConfirmMap(val, val2 string, b bool) {
	if _, ok := p.prePareConfirmCount[val]; !ok {
		p.prePareConfirmCount[val] = make(map[string]bool)
	}
	p.prePareConfirmCount[val][val2] = b
}

//为多重映射开辟赋值
func (p *pbft) setCommitConfirmMap(val, val2 string, b bool) {
	if _, ok := p.commitConfirmCount[val]; !ok {
		p.commitConfirmCount[val] = make(map[string]bool)
	}
	p.commitConfirmCount[val][val2] = b
}

//传入节点编号， 获取对应的公钥
func (p *pbft) getPubKey(nodeID string) []byte {
	key, err := ioutil.ReadFile("Keys/" + nodeID + "/" + nodeID + "_RSA_PUB")
	if err != nil {
		log.Panic(err)
	}
	return key
}

//传入节点编号， 获取对应的私钥
func (p *pbft) getPivKey(nodeID string) []byte {
	key, err := ioutil.ReadFile("Keys/" + nodeID + "/" + nodeID + "_RSA_PIV")
	if err != nil {
		log.Panic(err)
	}
	return key
}

// rsa 签名加密

//如果当前目录下不存在目录Keys，则创建目录，并为各个节点生成rsa公私钥
func genRsaKeys() {
	if !isExist("./Keys") {
		fmt.Println("检测到还未生成公私钥目录，正在生成公私钥 ...")
		err := os.Mkdir("Keys", 0644)
		if err != nil {
			log.Panic()
		}
		for i := 0; i <= nodeCount; i++ {
			if !isExist("./Keys/N" + strconv.Itoa(i)) {
				err := os.Mkdir("./Keys/N"+strconv.Itoa(i), 0644)
				if err != nil {
					log.Panic()
				}
			}
			priv, pub := getKeyPair()
			privFileName := "Keys/N" + strconv.Itoa(i) + "/N" + strconv.Itoa(i) + "_RSA_PIV"
			file, err := os.OpenFile(privFileName, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Panic(err)
			}
			defer file.Close()
			file.Write(priv)

			pubFileName := "Keys/N" + strconv.Itoa(i) + "/N" + strconv.Itoa(i) + "_RSA_PUB"
			file2, err := os.OpenFile(pubFileName, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Panic(err)
			}
			defer file2.Close()
			file2.Write(pub)
		}
		fmt.Println("已为节点们生成RSA公私钥")
	}
}

//生成rsa公私钥
func getKeyPair() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

//判断文件或文件夹是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

//数字签名
func (p *pbft) RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//签名验证
func (p *pbft) RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// tcp 用做监听
//客户端使用的tcp监听
func clientTcpListen() {
	listen, err := net.Listen("tcp", clientAddr)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(b))
	}

}

//节点使用的tcp监听
func (p *pbft) tcpListen() {
	listen, err := net.Listen("tcp", p.node.addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("节点开启监听，地址：%s\n", p.node.addr)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		p.handleRequest(b)
	}

}

//使用tcp发送消息
func tcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}

	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}

// 根据真实网络环境生成随机延迟
func RandomDelayGenerator() time.Duration {
	mathRand.Seed(time.Now().UnixNano())

	return time.Duration((80+mathRand.Int63n(40)) * int64(time.Millisecond))
}

func ReadFromCSV() map[string]string {
	b, err := os.Open("E:\\code1\\consensus\\pbft2\\networks4.csv")
	if err != nil {
		fmt.Print(err)
	}
	r1 := csv.NewReader(b)
	content, err := r1.ReadAll()
	if err != nil {
		log.Fatalf("can not readall, err is %+v", err)
	}
	for _, row := range content {
		fmt.Println(row)
	}

	networkMap := make(map[string]string)
	for i:=0; i<len(content); i++ {
		networkMap[content[i][0]] = content[i][1]
	}
	fmt.Println(networkMap)
	return networkMap
}


