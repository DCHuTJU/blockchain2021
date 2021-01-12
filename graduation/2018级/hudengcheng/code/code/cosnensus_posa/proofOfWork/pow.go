package proofOfWork

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"
)

//用于存放区块,以便连接成区块链
var blockchain []block

//挖矿难度值
var diffNum uint = 23

var totalTime []int64

type block struct {
	//上一个区块的Hash
	Lasthash string
	//本区块Hash
	Hash string
	//区块存储的数据（比如比特币UTXO模型 则此处可用于存储交易）
	Data string
	//时间戳
	Timestamp string
	//区块高度
	Height int
	//难度值
	DiffNum uint
	//随机数
	Nonce int64
}

//区块挖矿（通过自身递增nonce值计算hash）
func mine(data string) block {
	start := time.Now().UnixNano()
	if len(blockchain) < 1 {
		log.Panic("还未生成创世区块！")
	}
	lastBlock := blockchain[len(blockchain)-1]
	//制造一个新的区块
	newBlock := new(block)
	newBlock.Lasthash = lastBlock.Hash
	newBlock.Timestamp = time.Now().String()
	newBlock.Height = lastBlock.Height + 1
	newBlock.DiffNum = diffNum
	newBlock.Data = data
	var nonce int64 = 0
	//根据挖矿难度值计算的一个大数
	newBigint := big.NewInt(1)
	newBigint.Lsh(newBigint, 256-diffNum) //相当于左移 1<<256-diffNum
	for {
		newBlock.Nonce = nonce
		newBlock.getHash()
		hashInt := big.Int{}
		hashBytes, _ := hex.DecodeString(newBlock.Hash)
		hashInt.SetBytes(hashBytes) //把本区块hash值转换为一串数字
		//如果hash小于挖矿难度值计算的一个大数，则代表挖矿成功
		if hashInt.Cmp(newBigint) == -1 {
			break
		} else {
			nonce++ //不满足条件，则不断递增随机数，直到本区块的散列值小于指定的大数
		}
	}
	end := time.Now().UnixNano()
	fmt.Println("Time cost is: ", end - start)
	totalTime = append(totalTime, end - start)
	return *newBlock
}

func (b *block) serialize() []byte {
	bytes, err := json.Marshal(b)
	if err != nil {
		log.Panic(err)
	}
	return bytes
}

func (b *block) getHash() {
	result := sha256.Sum256(b.serialize())
	b.Hash = hex.EncodeToString(result[:])
}

func calculateAverage() float64 {
	sum := 0.0
	for i:=0; i<len(totalTime); i++ {
		sum += float64(totalTime[i])
	}

	return sum / float64(len(totalTime))
}
