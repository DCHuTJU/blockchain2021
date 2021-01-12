package merkleTree

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/cbergoon/merkletree"
)

type Content struct {
	Loc string
	// 单位为 MB
	Size string
	Dig string
	Hash string
}

// 计算交易 Hash
func (c Content) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(c.Loc + c.Size + c.Dig + c.Hash)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (c Content) Equals(other merkletree.Content) (bool, error) {
	return c.Loc == other.(Content).Loc && c.Size == other.(Content).Size && c.Dig == other.(Content).Dig && c.Hash == other.(Content).Hash, nil
}

func GenerateTransactionHash(loc, size, dig string) string {
	hash := sha256.New()
	str := loc + size + dig
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}