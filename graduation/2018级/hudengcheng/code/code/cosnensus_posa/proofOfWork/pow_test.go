package proofOfWork

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMine(t *testing.T) {
	//制造一个创世区块
	genesisBlock := new(block)
	genesisBlock.Timestamp = time.Now().String()
	genesisBlock.Lasthash = "0000000000000000000000000000000000000000000000000000000000000000"
	genesisBlock.Height = 1
	genesisBlock.Nonce = 0
	genesisBlock.DiffNum = 0
	genesisBlock.getHash()
	fmt.Println("The genesisBlock is: ", *genesisBlock)
	//将创世区块添加进区块链
	blockchain = append(blockchain, *genesisBlock)
	for i := 0; i < 10; i++ {
		newBlock := mine(strconv.Itoa(i))
		blockchain = append(blockchain, newBlock)
		fmt.Println(newBlock)
	}
	fmt.Println("Average time cost is: ", calculateAverage())
}