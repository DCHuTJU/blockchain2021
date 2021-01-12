package main

import (
	"fmt"
	"log"
	"os/exec"
)

//执行python脚本
func CmdPythonSaveImageDpi() (str string, err error) {
	args := "RL-NSGA2.py"
	rlt, err := exec.Command("python", args).Output()
	if err != nil {
		return "", err
	}
	result := string(rlt)
	fmt.Println("result is: ", result)
	return result, nil
}

func main() {
	rlt, err := CmdPythonSaveImageDpi()
	if err != nil {
		print("error is: ", err.Error())
		return
	}
	fmt.Println("rlt is: ", rlt)
	log.Println("调用成功")
}
