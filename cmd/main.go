package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/naumovrus/go-tron-api/internal/models"
	"github.com/sirupsen/logrus"
)

/*
TODO: refactor code
*/

var payload = strings.NewReader("{\"detail\":true}")

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Printf("Start parse transactions")
	// load .env files
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error load .env files: %s", err.Error())
	}

	// parse first block
	req, err := http.NewRequest("POST", os.Getenv("NODE_URL"), payload)
	if err != nil {
		logrus.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Fatal(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Printf("client: could not read response body: %s\n", err)

	}

	block := models.Block{}

	err = json.Unmarshal(resBody, &block)
	if err != nil {
		logrus.Printf("%s \n", err)
	}

	// set

	set := make(map[string]struct{})

	logrus.Printf("block id: %s number of block: %v", block.BlockID, block.BlockHeader.RawData.Number)
	logrus.Printf("block txTrieRoot: %s ", block.BlockHeader.RawData.TxTrieRoot)
	logrus.Printf("witness_address: %s", block.BlockHeader.RawData.WitnessAddress)
	logrus.Printf("witness_signature: %s", block.BlockHeader.WitnessSignature)
	for _, tx := range block.Transactions {
		set[tx.RawData.TransactionRawData[0].Type] = struct{}{}
	}

	sliceKeys := make([]string, 0, 1)

	for key := range set {
		sliceKeys = append(sliceKeys, key)
	}
	logrus.Printf("all types of transactions: %v count: %v", sliceKeys, len(block.Transactions))
	var nextBlock int32
	nextBlock = int32(block.BlockHeader.RawData.Number + 1)
	payloadGetBlockById := strings.NewReader(fmt.Sprintf("{\"num\":%d}", nextBlock))
	fmt.Println("--------------------------------------------------------")
	time.Sleep(time.Second * 5)

	go func() {

		for {
			set = make(map[string]struct{})
			payloadGetBlockById = strings.NewReader(fmt.Sprintf("{\"num\":%d}", nextBlock))
			logrus.Printf("get by num: %v", nextBlock)
			block = models.Block{}
			req, err := http.NewRequest("POST", os.Getenv("NODE_URL_NUM"), payloadGetBlockById)
			if err != nil {
				logrus.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logrus.Fatal(err)
			}

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				logrus.Printf("client: could not read response body: %s\n", err)

			}

			err = json.Unmarshal(resBody, &block)
			if err != nil {
				logrus.Printf("%s \n", err)
			}

			log.Printf("block id: %s\nnumber of block: %v \n", block.BlockID, block.BlockHeader.RawData.Number)
			for _, tx := range block.Transactions {
				set[tx.RawData.TransactionRawData[0].Type] = struct{}{}
			}

			sliceKeys := make([]string, 0, 1)

			for key := range set {
				sliceKeys = append(sliceKeys, key)
			}
			logrus.Printf("block txTrieRoot: %s ", block.BlockHeader.RawData.TxTrieRoot)
			logrus.Printf("all types of transactions: %v count: %v", sliceKeys, len(block.Transactions))
			logrus.Printf("witness_address: %s", block.BlockHeader.RawData.WitnessAddress)
			logrus.Printf("witness_signature: %s", block.BlockHeader.WitnessSignature)
			fmt.Println("--------------------------------------------------------")
			nextBlock += 1

			time.Sleep(time.Second * 5)

		}

	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
