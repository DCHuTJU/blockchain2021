package state

// 区块结构
type Block struct {
	Height int    `json:"height"`
	Data   string `json:"data"`
	MerkleRoot string `json:"merkleroot"`
	AllocationResult string `json:"allocationresult"`
	BlockBody []Transaction
}

// 交易内容定义
type Transaction struct {

}

// MerkleRoot 构建过程


