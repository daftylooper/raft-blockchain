package main

import (
	"encoding/json"
	"fmt"
	"hash"
	"time"
)

type transaction struct {
	Index     int32
	Timestamp time.Time
	Sender    string
	Receiver  string
	Amount    float32
}

type block struct {
	Index        int32
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Transactions []transaction
	Nonce        int32
}

type chain struct {
	Height     int32
	Difficulty int32
	Blocks     []block
}

func NewTransaction(Index int32, Sender string, Receiver string, Amount float32) *transaction {
	tx := new(transaction)
	tx.Index = Index
	tx.Sender = Sender
	tx.Receiver = Receiver
	tx.Amount = Amount
	tx.Timestamp = time.Now()

	return tx
}

func NewBlock(Index int32, PreviousHash string, Nonce int32) *block {
	blk := new(block)
	blk.Index = Index
	//add timestamp once a block is manufactured
	blk.PreviousHash = PreviousHash
	blk.Nonce = Nonce

	return blk
}

// method on block to add a transaction
func (Block *block) AppendTransaction(Sender string, Receiver string, Amount float32) {
	NewIndex := int32(len(Block.Transactions))
	newTX := NewTransaction(NewIndex, Sender, Receiver, Amount)
	Block.Transactions = append(Block.Transactions, *newTX)
}

func (Block *block) ForgeBlock(h hash.Hash, Difficulty int32) {
	//generate hash by string(block params...)
	h.Write(Block.MarshalBlock())
	BlockHash := string(h.Sum(nil))
	//fmt.Printf("%x\n", BlockHash)
	Block.Hash = BlockHash
	//bruteforce the nonce
	POWBlockHash := BlockHash
	Nonce := 0
	for !DifficultySatisfied(POWBlockHash, Difficulty) {
		h.Write([]byte(POWBlockHash + fmt.Sprintf("%q", Nonce)))
		POWBlockHash = string(h.Sum(nil))
		Nonce += 1
	}
	// t := fmt.Sprintf("%x", POWBlockHash)
	// fmt.Println(t)
	Block.Nonce = int32(Nonce)
	//fmt.Println(Nonce)
	Block.Timestamp = time.Now()
	//append block to the chain

}

func DifficultySatisfied(CurrentHash string, Difficulty int32) bool {
	for i := 0; i < int(Difficulty); i++ {
		if CurrentHash[i] != 0 {
			return false
		}
	}
	return true
}

// marshall all fields of a block and return a string
func (Block *block) MarshalBlock() []byte {
	marshalledTransaction, txmerr := json.Marshal(Block.Transactions)
	if txmerr != nil {
		fmt.Println("Marshall Error:", txmerr)
	}
	marshalledTimestamp, tsmerr := json.Marshal(Block.Timestamp)
	if tsmerr != nil {
		fmt.Println("Marshall Error:", tsmerr)
	}
	marshalled := string(Block.Index) + string(Block.PreviousHash) + string(marshalledTimestamp) + string(marshalledTransaction)
	//fmt.Println(marshalled)

	return []byte(marshalled)
}

func NewChain(Difficulty int32) *chain {
	ch := new(chain)
	ch.Difficulty = Difficulty
	return ch
}

func (Chain *chain) AppendBlock(Block *block) {
	Chain.Blocks = append(Chain.Blocks, *Block)
	Chain.Height += 1
}

func (Chain *chain) GetPreviousHash() string {
	if Chain.Height < 1 {
		return "0"
	}
	return Chain.Blocks[Chain.Height-1].Hash
}

func (Chain *chain) DebugChain() {
	fmt.Println("Chain Height -", Chain.Height)
	fmt.Println("Chain Difficulty -", Chain.Difficulty)
	fmt.Println("Chain Blocks -")
	for i := 0; i < int(Chain.Height); i++ {
		HexHash := fmt.Sprintf("%x", Chain.Blocks[i].Hash)
		PrevHexHash := fmt.Sprintf("%x", Chain.Blocks[i].PreviousHash)
		fmt.Println("\nBLOCK", i, ":\n--------\n", HexHash, Chain.Blocks[i].Index, Chain.Blocks[i].Nonce, PrevHexHash, Chain.Blocks[i].Timestamp, Chain.Blocks[i].Transactions)
	}
}

func (Chain *chain) CheckIntegrity() (int, bool) {
	for i := 1; i < int(Chain.Height); i++ {
		if Chain.Blocks[i].PreviousHash != Chain.Blocks[i-1].Hash {
			return i, false
		}
	}
	return -1, true
}
