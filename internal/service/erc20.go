package service

import (
  "context"
  "log"
  "math/big"
  "strings"

  "github.com/ethereum/go-ethereum"
  "github.com/ethereum/go-ethereum/accounts/abi"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/ethclient"
)

func GetBalanceOf(account, contract string, client *ethclient.Client) *big.Int {
  const erc20ABI = `[{
    "constant": true,
    "inputs": [{"name":"owner","type":"address"}],
    "name": "balanceOf",
    "outputs": [{"name":"","type":"uint256"}],
    "type": "function"
  }]`
  ctx := context.Background()

  tokenAddr := common.HexToAddress(contract)
  userAddr := common.HexToAddress(account)

  // 解析 ABI
  parsed, err := abi.JSON(strings.NewReader(erc20ABI))
  if err != nil {
    log.Fatal(err)
  }

  // 打包 call data
  data, err := parsed.Pack("balanceOf", userAddr)
  if err != nil {
    log.Fatal(err)
  }

  // 构造 eth_call 消息
  msg := ethereum.CallMsg{
    To:   &tokenAddr,
    Data: data,
  }

  // 执行 eth_call（不会消耗 gas）
  result, err := client.CallContract(ctx, msg, nil)
  if err != nil {
    log.Fatal(err)
  }

  // 解析返回值
  return new(big.Int).SetBytes(result)
}
