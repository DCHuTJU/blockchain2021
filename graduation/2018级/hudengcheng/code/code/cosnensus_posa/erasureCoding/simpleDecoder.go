package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/klauspost/reedsolomon"
)

var dataShardsDecoder = flag.Int("data", 4, "Number of shards to split the data into")
var parShardsDecoder = flag.Int("par", 2, "Number of parity shards")
var outFile = flag.String("out", "", "Alternative output path/file")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  simple-decoder [-flags] basefile.ext\nDo not add the number to the filename.\n")
		fmt.Fprintf(os.Stderr, "Valid flags:\n")
		flag.PrintDefaults()
	}
}

func main() {
	// Parse flags
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: No filenames given\n")
		flag.Usage()
		os.Exit(1)
	}
	fname := args[0]
	start := time.Now().UnixNano()
	// Create matrix
	enc, err := reedsolomon.New(*dataShardsDecoder, *parShardsDecoder)
	checkErrDecoder(err)

	// Create shards and load the data.
	shards := make([][]byte, *dataShardsDecoder+*parShardsDecoder)
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", fname, i)
		fmt.Println("Opening", infn)
		shards[i], err = ioutil.ReadFile(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
		}
	}

	// Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		err = enc.Reconstruct(shards)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			os.Exit(1)
		}
		checkErrDecoder(err)
	}

	// Join the shards and write them
	outfn := *outFile
	if outfn == "" {
		outfn = fname
	}

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	checkErrDecoder(err)

	// We don't know the exact filesize.
	err = enc.Join(f, shards, len(shards[0])**dataShardsDecoder)
	checkErrDecoder(err)

	end := time.Now().UnixNano()
	fmt.Println("Time cost is: ", end - start)
}

func checkErrDecoder(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}

