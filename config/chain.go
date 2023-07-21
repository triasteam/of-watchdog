package config

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/openfaas/of-watchdog/logger"
)

type Chain struct {
	Id                 int64
	Addr               string
	FunctionClientAddr string
	FunctionOracleAddr string
	KeyFilePath        string
	KeyPassword        string
}

func LoadChainConfig() *Chain {
	return &Chain{}
}

func (c Chain) ChainID() int64 {
	return c.Id
}

func (c Chain) ChainAddr() string {
	return c.Addr
}

func (c Chain) Key() *keystore.Key {

	keyBytes, err := os.ReadFile(c.KeyFilePath)
	if err != nil {
		logger.Fatal("failed to read key file", "err", err)
	}

	key, err := keystore.DecryptKey(keyBytes, c.KeyPassword)
	if err != nil {
		logger.Fatal("wrong key", "key", string(keyBytes), "pw", c.KeyPassword)
	}
	return key
}

func (c Chain) FuncClientAddr() string {
	return c.FunctionClientAddr
}

func (c Chain) FuncOracleClientAddr() string {
	return c.FunctionOracleAddr
}
