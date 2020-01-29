package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/INFURA/go-ethlibs/eth"
)

type blockAndTx struct {
	blockHash string
	txHash    []string
}

func main() {

	draw := flag.Int("draw", 5, "number of blocks drawn")
	baseURL := flag.String("url", "http://localhost:8000", "url of the API")
	stdout := flag.Bool("stdout", false, "Write targets to stdout")

	flag.Parse()

	if !*stdout {
		defer timeTrack(time.Now(), "OVERALL")
	}
	lastURL := *baseURL + "/block/last/height"
	resp, err := http.Get(lastURL)

	check(err)

	defer resp.Body.Close()

	var result map[string]int64
	json.NewDecoder(resp.Body).Decode(&result)

	lastHeight := result["lastBlockHeight"]

	var hashes = make(chan blockAndTx)
	var blockHashes []string
	var blockIds []int64
	var transactionHashes []string
	for i := 0; i <= *draw; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(lastHeight))
		check(err)

		go getTxfromBlock(*stdout, i, n.Int64(), hashes)
		blockIds = append(blockIds, n.Int64())
	}
	for i := 0; i < *draw; i++ {
		select {
		case h := <-hashes:
			blockHashes = append(blockHashes, h.blockHash)
			transactionHashes = append(transactionHashes, h.txHash...)
		}
	}
	if !*stdout {
		fmt.Printf("We gathered %d block hash and %d tx hash \n", len(blockHashes), len(transactionHashes))
		fmt.Println("printing targets in targets.txt")
	}
	writeTargetFile(*stdout, *baseURL, blockHashes, blockIds, transactionHashes)

}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeTargetFile(stdout bool, baseURL string, blockHashes []string, blockIds []int64, transactionHashes []string) {
	targets := "GET " + baseURL + "/block/last/height\n"

	for _, b := range blockHashes {
		targets += fmt.Sprintf("GET %s/block/%s\n", baseURL, b)
	}
	for _, t := range transactionHashes {
		targets += fmt.Sprintf("GET %s/transaction/%s\n", baseURL, t)
	}
	for _, bi := range blockIds {
		targets += fmt.Sprintf("GET %s/block/%d\n", baseURL, bi)
	}

	if stdout {
		fmt.Print(targets)
	} else {
		d1 := []byte(targets)
		err := ioutil.WriteFile("./targets.txt", d1, 0644)
		check(err)
	}
}
func getTxfromBlock(stdout bool, index int, height int64, hashes chan<- blockAndTx) {
	if !stdout {
		defer timeTrack(time.Now(), "getTxfromBlock")
	}

	urlToBlock := fmt.Sprintf("http://localhost:8000/block/%d", height)
	resp, err := http.Get(urlToBlock)
	if err != nil {
		log.Println("***************************", err)
	}
	var block eth.Block
	txHash := []string{}
	json.NewDecoder(resp.Body).Decode(&block)
	for id, t := range block.Transactions {
		if !stdout {
			fmt.Printf("	id:%d txhash:%+v\n", id, t.Transaction.Hash)
		}
		txHash = append(txHash, t.Transaction.Hash.String())
	}
	res := blockAndTx{
		blockHash: block.Hash.String(),
		txHash:    txHash,
	}
	hashes <- res
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
