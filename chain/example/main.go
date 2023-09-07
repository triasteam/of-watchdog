package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/openfaas/of-watchdog/chain/actor"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/openfaas/of-watchdog/config"
	"github.com/openfaas/of-watchdog/logger"
)

func main() {
	funcNameBytes, err := hex.DecodeString("1230000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		logger.Error("failed to decode function name")
		return
	}
	funcName := [32]byte(funcNameBytes)
	fmt.Println(funcName)
	source := "sdfssf"
	secrets := []byte{1, 2, 3}
	args := []string{"test"}

	config.SetEnv(
		12345678, "ws://210.73.218.170:9546",
		"0xE8935Af625542fc62A866784F5f490636d9dbC66", "0x0000000000000000000000000000000000002004",
		"", "test",
		"./keystore/UTC--2023-06-05T09-50-10.886531000Z--989777e983d4fccba32d857d797fdb75c27571c5", "123456",
	)

	chainConfig := config.LoadChainConfig()

	ethCli, err := ethclient.Dial(chainConfig.Addr)
	if err != nil {
		logger.Error("failed to connect to ", "node addr", chainConfig.Addr)
		return
	}

	funcCli, err := NewFunctionConsumer(common.HexToAddress(chainConfig.FuncClientAddr()), ethCli)
	if err != nil {
		logger.Error("failed to new function consumer", "err", err)
		return
	}

	oracleCli, err := actor.NewFunctionOracle(common.HexToAddress(chainConfig.FuncOracleClientAddr()), ethCli)
	if err != nil {
		logger.Error("failed to new function consumer", "err", err)
		return
	}

	go WatchEvent(ethCli, oracleCli, funcCli, chainConfig)

	auth, err := bind.NewKeyedTransactorWithChainID(chainConfig.Key().PrivateKey, new(big.Int).SetInt64(chainConfig.ChainID()))
	if err != nil {
		logger.Error("failed to bind tx", "err", err)
		return
	}

	sink := make(chan *FunctionConsumerRequestSent)
	defer close(sink)
	respSub, err := funcCli.WatchRequestSent(&bind.WatchOpts{Context: context.Background()}, sink, nil, nil)
	if err != nil {
		logger.Error("failed to watch RequestSent event", "err", err)
		return
	}
	defer respSub.Unsubscribe()
	logger.Info("start to send tx")
	tx, err := funcCli.ExecuteRequest(&bind.TransactOpts{
		From:   auth.From,
		Signer: auth.Signer,
	}, funcName, source, secrets, args)
	if err != nil {
		logger.Error("failed to send execute request tx", "err", err)
		return
	}
	var (
		reqSent *FunctionConsumerRequestSent
	)
	logger.Info("select event")
	select {
	case reqSent = <-sink:
	case err = <-respSub.Err():
		logger.Error("failed to send exec request tx", "err", err)
		return
	}
	logger.Info("fulfilled request", "tx hash", tx.Hash().String(), "blockNumber", reqSent.Raw.BlockNumber, "raw data", hex.Dump(reqSent.Raw.Data))
	for {
		time.Sleep(time.Second * 10)
	}
}
