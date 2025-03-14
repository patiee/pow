package main

import (
	"fmt"
	"time"

	"github.com/patiee/pow/model"
)

func main() {
	// Test creating a genesis block
	genesis := model.Block{
		Height:       0,
		Timestamp:    time.Now().Unix(),
		PreviousHash: "0",
		Nonce:        0,
		Difficulty:   "1",
	}
	fmt.Printf("Starting to mine block %+v\n", genesis)
	start := time.Now().Unix()
	genesisHash, err := genesis.MineBlock()
	if err != nil {
		fmt.Printf("Error mining block: %s\n", err)
		return
	}
	fmt.Printf("Genesis Block: %+v\n", genesis)
	fmt.Printf("Hash: %s\n", genesisHash)
	fmt.Printf("Block minted, took %d \n", time.Now().Unix()-start)
}
