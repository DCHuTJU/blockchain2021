package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/klauspost/reedsolomon"
)

var dataShards = flag.Int("data", 4, "Number of shards to split the data into, must be below 257.")
var parShards = flag.Int("par", 2, "Number of parity shards")
var outDir = flag.String("out", "", "Alternative output directory")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  simple-encoder [-flags] filename.ext\n\n")
		fmt.Fprintf(os.Stderr, "Valid flags:\n")
		flag.PrintDefaults()
	}
}

func main() {
	// Parse command line parameters.
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: No input filename given\n")
		flag.Usage()
		os.Exit(1)
	}
	if (*dataShards + *parShards) > 256 {
		fmt.Fprintf(os.Stderr, "Error: sum of data and parity shards cannot exceed 256\n")
		os.Exit(1)
	}
	fname := args[0]

	start := time.Now().UnixNano()
	// Create encoding matrix.
	enc, err := reedsolomon.New(*dataShards, *parShards)
	checkErr(err)

	fmt.Println("Opening", fname)
	b, err := ioutil.ReadFile(fname)
	checkErr(err)

	// Split the file into equally sized shards.
	shards, err := enc.Split(b)
	checkErr(err)
	fmt.Printf("File split into %d data+parity shards with %d bytes/shard.\n", len(shards), len(shards[0]))

	// Encode parity
	err = enc.Encode(shards)
	checkErr(err)

	// Write out the resulting files.
	dir, file := filepath.Split(fname)
	if *outDir != "" {
		dir = *outDir
	}
	for i, shard := range shards {
		outfn := fmt.Sprintf("%s.%d", file, i)

		fmt.Println("Writing to", outfn)
		err = ioutil.WriteFile(filepath.Join(dir, outfn), shard, 0644)
		checkErr(err)
	}

	end := time.Now().UnixNano()
	fmt.Println("Time cost is: ", end - start)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
