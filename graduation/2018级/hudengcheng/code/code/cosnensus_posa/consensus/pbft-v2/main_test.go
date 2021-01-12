package main

import (
	"testing"
)

func TestReadFromCSV(t *testing.T) {
	ReadFromCSV()
}

func TestJointMessage(t *testing.T) {
	tmp := jointMessage(cRequestCommit, []byte(""))
	t.Log(string(tmp))
}

func TestSplitMessage(t *testing.T) {
	cmd, content := splitMessage([]byte("requestcommit  "))
	t.Log("cmd is: ", cmd)
	t.Log("content is: ", content)
}

func TestGetDigest(t *testing.T) {
	tmp := Request{

	}
	rlt := getDigest(tmp)
	t.Log(rlt)
	t.Log(len(rlt))
}