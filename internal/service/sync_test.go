package service

import (
  "fmt"
  "testing"

  "github.com/ethereum/go-ethereum/crypto"
)

func TestEventHash(t *testing.T) {
  transfer := []byte("Transfer(address,address,uint256)")
  approval := []byte("Approval(address,address,uint256)")
  mint := []byte("Mint(address,uint256)")
  burn := []byte("Burn(address,uint256)")
  transferHash := crypto.Keccak256Hash(transfer)
  approvalHash := crypto.Keccak256Hash(approval)
  mintHash := crypto.Keccak256Hash(mint)
  burnHash := crypto.Keccak256Hash(burn)
  fmt.Println("approval hash:", approvalHash.Hex()) // 对应 topics[0]
  fmt.Println("transfer hash:", transferHash.Hex()) // 对应 topics[0]
  fmt.Println("mint hash:", mintHash.Hex())         // 对应 topics[0]
  fmt.Println("burn hash:", burnHash.Hex())         // 对应 topics[0]
}
