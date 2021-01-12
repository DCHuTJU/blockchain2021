package util

import (
	"log"
	"os"
)

const(
	// 文件大小 GB
	fileSize = 0.3
)

// 300 MB
var size = int64(1024 * 1024 * 1025 * fileSize)

// 生成文件工具，需要自己制定 生成文件大小
func ExampleTruncate() {
	f, err := os.Create("foobar1.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := f.Truncate(size); err != nil {
		log.Fatal(err)
	}
	// Output:
	//
}

func ExampleSeek() {
	f, err := os.Create("foobar2.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Seek(size-1, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte{0})
	if err != nil {
		log.Fatal(err)
	}
}
