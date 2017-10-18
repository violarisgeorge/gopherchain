package main

import (
	"crypto/sha256"
	"time"
	"encoding/hex"
	"strconv"
	"container/list"
	"fmt"
	"os"
  "math/rand"
)

type Block struct {

	index int
	timestamp time.Time
	data string
	previousHash string

}

func genesis() Block {

	var sha = sha256.New()
	sha.Write([]byte("This is Gopher, building your chain!"))

	genesisBlock := Block {0, time.Now(), "This is the Genesis block", hex.EncodeToString(sha.Sum(nil))}

	return genesisBlock

}

func nextBlock(lastBlock Block) Block {

	blockIndex := lastBlock.index + 1
	blockTime := time.Now().String()
	blockData := "This is block " + strconv.Itoa(blockIndex)
	previousHash := lastBlock.previousHash

	var sha = sha256.New()
	tryString := lastBlock.previousHash
	rand.Seed(time.Now().UnixNano())
	nonce := make([]byte, 4)
	rand.Read(nonce)

	hashString := "00000" + lastBlock.previousHash[5:len(lastBlock.previousHash)]

	for ((tryString[:5]) != hashString[:5]) {


			sha.Write([]byte(tryString))
			tryString = hex.EncodeToString(sha.Sum(nil))
			tryString = tryString + hex.EncodeToString(nonce)
			fmt.Printf(tryString)
			fmt.Printf("\n\n")

			nonce[3]++
	}

 	blockString := strconv.Itoa(blockIndex) + blockTime + blockData + previousHash + tryString
	var sha2 = sha256.New()
	sha2.Write([]byte(blockString))

	t, err := os.Create("height")
	check(err)
	n, err := t.WriteString(strconv.Itoa(blockIndex))
	fmt.Printf("Written height to file", n)
	t.Sync()

	return Block {blockIndex, time.Now(), blockData, hex.EncodeToString(sha2.Sum(nil)) }

}


func check(e error) {
	if e != nil {
			panic(e)
	}
}

func main() {

	var blockchain = list.New()
	var genesisBlock = genesis()
	blockchain.PushBack(genesisBlock)
	var previousBlock = genesisBlock

	f, err := os.OpenFile("dat", os.O_APPEND|os.O_WRONLY|
os.O_CREATE, 0600)
	check(err)

	for e:= blockchain.Front(); e != nil; e = e.Next() {

		newBlock := nextBlock(previousBlock)

		blockchain.PushBack(newBlock)

		fmt.Printf("[HEIGHT]: ")
		fmt.Printf(strconv.Itoa(newBlock.index))
		fmt.Printf("\n")
		fmt.Printf("Block ")
		fmt.Printf(newBlock.previousHash)
		fmt.Printf(" has been added to the blockchain!\n")

		heightLabel := "[HEIGHT]: "
		height := strconv.Itoa(newBlock.index)
		blockLabel := "Block: "
		block := newBlock.previousHash

		blockInfo := heightLabel + height + " \n" + blockLabel + block + "\n"
		n, err := f.WriteString(blockInfo)
		check(err)
		fmt.Printf("wrote %d bytes\n", n)
		f.Sync()

		previousBlock = newBlock

		time.Sleep(10000 * time.Millisecond) //simulate block creation by delaying output

	}

}
