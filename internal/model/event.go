package model

import (
  "math/big"

  "github.com/ethereum/go-ethereum/common"
)

type TransferEvent struct {
  From  common.Address
  To    common.Address
  Value *big.Int
}

type ApprovalEvent struct {
  Owner   common.Address
  Spender common.Address
  Value   *big.Int
}

type MintEvent struct {
  Addr   common.Address
  Amount *big.Int
}

type BurnEvent struct {
  Addr   common.Address
  Amount *big.Int
}
