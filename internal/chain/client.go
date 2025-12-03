package chain

import (
  "errors"

  "com.wh1200.points/internal/config"
  "github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
  Client *ethclient.Client
}

const (
  EthChainID     = 1
  SepoliaChainId = 11155111
)

func New(chainId uint64, nodeUrl string) (*Client, error) {
  switch chainId {
  case EthChainID, SepoliaChainId:
    client, err := ethclient.Dial(nodeUrl + config.Cfg.NodeApiKey)
    if err != nil {
      return nil, err
    }
    return &Client{
      Client: client,
    }, nil
  default:
    return nil, errors.New("unsupported chain id")
  }
}
