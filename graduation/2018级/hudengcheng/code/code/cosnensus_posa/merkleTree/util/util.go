package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
)

type FileChunk struct {
	Loc string
	MD5 string
	Size string
	IsPar string
}

// 大文件 md5 计算片段大小
const Filechunk = 8192

func GenerateLocHash(loc string) string {
	hash := sha256.New()
	hash.Write([]byte(loc))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

// 小文件 md5
func GenerateDigHashL(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}

	md5Str := fmt.Sprintf("%x", md5hash.Sum(nil))
	return md5Str
}

// 大文件 md5
func GenerateDigHashB(filepath string) string {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(Filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(Filechunk, float64(filesize-int64(i*Filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return string(hash.Sum(nil))
}

func CalculateMD5(fileList []FileChunk) string {

	strList := []string{}
	MD5List := []string{}

	for i:=0; i<len(fileList); i++ {
		str := fileList[i].Loc + fileList[i].MD5 + fileList[i].Size + fileList[i].IsPar
		strList = append(strList, str)
	}

	for i:=0; i<len(strList); i++ {
		MD5 := md5.New()
		MD5.Write([]byte(strList[i]))
		bytes := MD5.Sum(nil)
		hashCode := hex.EncodeToString(bytes)
		MD5List = append(MD5List, hashCode)
	}

	tmp := ""
	for i:=0; i<len(MD5List); i++ {
		tmp += MD5List[i]
	}

	MD5 := md5.New()
	MD5.Write([]byte(tmp))
	bytes := MD5.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

func CalculateSHA256(str string) string {
	sha := sha256.New()
	sha.Write([]byte(str))
	bytes := sha.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}