package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	h := sha256.New()
	chain := NewChain(2)
	fmt.Println("init...")
	block := NewBlock(chain.Height, chain.GetPreviousHash(), 123)
	block.AppendTransaction("A", "B", 0.5)
	block.AppendTransaction("B", "A", 0.1)
	block.AppendTransaction("A", "C", 0.7)
	block.ForgeBlock(h, chain.Difficulty)
	//fmt.Println(block.Index, block.Timestamp, block.Nonce, block.Transactions)
	chain.AppendBlock(block)

	block = NewBlock(chain.Height, chain.GetPreviousHash(), 123)
	block.AppendTransaction("Q", "E", 0.411)
	block.AppendTransaction("B", "R", 0.10101)
	block.AppendTransaction("J", "C", 0.6999)
	block.ForgeBlock(h, chain.Difficulty)
	chain.AppendBlock(block)

	chain.DebugChain()
	// chain.Blocks[1].PreviousHash = "ec16b226d06747b0274d2ac22a28ec6bc390267323a588a5f588da06f69ba8ed"

	where, safe := chain.CheckIntegrity()

	if safe {
		fmt.Println("Chain is SAFE.")
	} else {
		fmt.Println("Chain Compromised. Where:", where)
	}
}
