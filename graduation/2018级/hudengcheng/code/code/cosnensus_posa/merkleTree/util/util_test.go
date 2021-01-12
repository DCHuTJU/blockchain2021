package util

import (
	"fmt"
	"testing"
)

func TestGenerateDigHash(t *testing.T) {
	fmt.Println(GenerateDigHashL("E:\\consensus\\merkleTree\\run.go"))
	fmt.Println(GenerateDigHashB("E:\\consensus\\erasureCoding\\foobar2.bin"))
}