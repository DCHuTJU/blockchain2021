package main

import (
	"bufio"
	"code/proofOfWork/tpsBlockchain/core"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"runtime"

	_ "net/http/pprof"
	"os"
	"time"
)

//var address = flag.String("ip", fmt.Sprintf("%s:%s", core.GetIpAddress()[0], core.BLOCKCHAIN_PORT), "")
var address = flag.String("ip", fmt.Sprintf("%s:%s", "127.0.0.1", core.BLOCKCHAIN_PORT), "")


func init() {
	flag.Parse()
}

func main() {
	// runtime.GOMAXPROCS(1)
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	core.Start(*address)
	//ReadStdin
	/*
		for {
			str := <-ReadStdin()
			core.Core.Blockchain.TransactionsQueue <- core.CreateTransaction(str)
		}
	*/
	tx := CreateTransactionTest("0.0001BTC")
	// if tx.VerifyTransaction(core.TRANSACTION_POW){
	// 	fmt.Println("Sig verify success!")
	// }
	go func() {
		for {
			for i := 0; i < core.TXPOOL_SIZE; i++ {
				core.Core.Blockchain.TransactionsQueue <- tx
			}
			fmt.Printf(".................................................pre-generating %d transactions........................................\n", core.TXPOOL_SIZE)
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// 定时监听
	go func() {
		for {
			v, _ := mem.VirtualMemory()
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
			fmt.Println(v)
		}
	}()

	for {
		<-ReadStdin()
	}
	// N := 500
	// txCh := RandomNTx(N)
	// for i := 0; i < N; i++ {
	// 	//time.Sleep(time.Microsecond * 30)
	// 	time.Sleep(time.Second * 1)
	// 	core.Core.Blockchain.TransactionsQueue <- <-txCh
	// }
	core.Wg.Wait()


}

func ReadStdin() chan string {

	cb := make(chan string)
	sc := bufio.NewScanner(os.Stdin)

	go func() {
		if sc.Scan() {
			cb <- sc.Text()
		}
	}()

	return cb
}

func RandomNTx(N int) chan *core.Transaction {
	txCh := make(chan *core.Transaction, N)
	for i := 0; i < N; i++ {
		signedTx := CreateTransactionTest("0.00001BTC")
		txCh <- signedTx
	}
	return txCh
}

func CreateTransactionTest(txt string) *core.Transaction {
	fromKey, toKey := core.GenerateNewKeypair(), core.GenerateNewKeypair()

	tx := core.NewTransaction(fromKey.Public, toKey.Public, []byte(txt))
	tx.Header.Nonce = tx.GenerateNonce(core.TRANSACTION_POW)

	tx.Signature = tx.Sign(fromKey)

	return tx
}
