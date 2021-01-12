package consensus

import "github.com/w3liu/consensus/types"

func Consensus() types.Result {
	reward, allocation := StorageAllocation()
	rlt := types.Result{
		Reward:     reward,
		Id:         "127.0.0.1:8003",
		Allocation: allocation,
	}
	return rlt
}

func StorageAllocation() (float64, []int) {
	return 0.001, []int{0, 1, 0, 1, 0, 0}
}