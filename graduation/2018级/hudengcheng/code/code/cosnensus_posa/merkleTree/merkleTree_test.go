package merkleTree

import (
	"code/merkleTree/util"
	"fmt"
	"github.com/cbergoon/merkletree"
	"log"
	"testing"
	"time"
	"unsafe"
)

func TestMerkleTree(t *testing.T) {
	filepath := "E:\\consensus\\erasureCoding\\foobar2.bin"
	str := util.GenerateDigHashL(filepath)
	var fileList []util.FileChunk
	start := time.Now()
	fmt.Println(start)
	fileList = append(fileList, util.FileChunk{
		Loc: "1",
		MD5: str,
		Size: "200",
		IsPar: "0",
	})

	fileList = append(fileList, util.FileChunk{
		Loc: "3",
		MD5: str,
		Size: "200",
		IsPar: "0",
	})

	fileList = append(fileList, util.FileChunk{
		Loc: "4",
		MD5: str,
		Size: "200",
		IsPar: "0",
	})

	fileList = append(fileList, util.FileChunk{
		Loc: "7",
		MD5: str,
		Size: "200",
		IsPar: "0",
	})

	fileList = append(fileList, util.FileChunk{
		Loc: "8",
		MD5: str,
		Size: "200",
		IsPar: "0",
	})

	fileList = append(fileList, util.FileChunk{
		Loc: "12",
		MD5: str,
		Size: "200",
		IsPar: "1",
	})

	fileList = append(fileList, util.FileChunk{
		Loc: "13",
		MD5: str,
		Size: "200",
		IsPar: "1",
	})

	var fileChunk util.FileChunk
	size := unsafe.Sizeof(fileChunk)
	fmt.Println("File chunk size is: ", size)

	var list []merkletree.Content
	//list = append(list, Content{
	//	Loc: "53bfc4015cf3e3be767498b5fbf4d1cad4802f46eff0c2830376b02b5317e3ad",
	//	Size: "200",
	//	// 用 md5 做的签名
	//	Dig: util.CalculateMD5(fileList),
	//	Hash: "ce922519a3c3ecaf9b0986c2449c7680895c15f4b0e9818e994e14a4d28b6aaf",
	//})

	for i:=0; i<40000; i++ {
		list = append(list, Content{
			Loc: util.CalculateSHA256(fmt.Sprintf("%d", i)),
			Size: "200",
			Dig: util.CalculateMD5(fileList),
			Hash: util.CalculateSHA256(util.CalculateMD5(fileList)+fmt.Sprintf("%d", i)),
		})
	}

	//list = append(list, Content{40
	//	Loc: "11672cafb21af3099f07ae32caf1aaf726cb72b996c0d244a6e3e491a952ef41",
	//	Size: "200",
	//	// 用 md5 做的签名
	//	Dig: CalculateMD5(fileList),
	//	Hash: "00e409723fe4d7cb130c27c881ba9dbc442396e9a5ec293d03557134604bfe77",
	//})
	//
	//list = append(list, Content{
	//	Loc: "17ccb6e0f49d1c941b994fa9b447dc86f967c125ff027dd8e411aa0a2a3a2473",
	//	Size: "200",
	//	// 用 md5 做的签名
	//	Dig: CalculateMD5(fileList),
	//	Hash: "892f703427073865e64d9a6842482c04d25c392edda396ffe1e825b5d02b311f",
	//})
	//
	//list = append(list, Content{
	//	Loc: "25295402fdd8de1d6cfd640be6c50b83ca148de80b564ce32b63c2c9f91d9346",
	//	Size: "200",
	//	// 用 md5 做的签名
	//	Dig: CalculateMD5(fileList),
	//	Hash: "b1681b1e8ec67dc18cc2524461f1593c2b92c70f0aee883edef141cd5c2d7bd5",
	//})
	//
	//list = append(list, Content{
	//	Loc: "448f3208af21d1138d4b58be9692c68a4c6ff37d57a9aed6f494a07f41d476a1",
	//	Size: "200",
	//	// 用 md5 做的签名
	//	Dig: CalculateMD5(fileList),
	//	Hash: "8536e420aa40d34956c1d46d1f55a0eb3ec70d22419b5fa9a6803a6b7f166f51",
	//})

	// 创建 Merkle Tree
	print("Length of MerkleTree is: ", len(list))
	tree, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}
	println("List size is: ", unsafe.Sizeof(t))
	end := time.Now()
	fmt.Println(end)
	fmt.Println("Time cost is: ", end.Sub(start))
	mr := tree.MerkleRoot()
	log.Println(mr)

	vt, err := tree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Verify Tree: ", vt)

	vc, err := tree.VerifyContent(list[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Verify Content: ", vc)

	// fmt.Println(t.Leafs)
}