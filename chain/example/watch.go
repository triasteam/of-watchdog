package main

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/openfaas/of-watchdog/chain"
	"github.com/openfaas/of-watchdog/chain/actor"
	"github.com/openfaas/of-watchdog/chain/cbor"
	"github.com/openfaas/of-watchdog/logger"
	"github.com/pkg/errors"
)

func WatchEvent(ethCli *ethclient.Client, oracleCli *actor.FunctionOracle, functionClient *FunctionConsumer, cfg chain.Configure) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(cfg.FuncClientAddr()),
			common.HexToAddress(cfg.FuncOracleClientAddr()),
		},
		Topics: [][]common.Hash{},
	}
	logs := make(chan types.Log)
	var (
		sub ethereum.Subscription
		err error
	)

	logger.Info("start subscribe")
	sub, err = ethCli.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logger.Error("failed to finish to subscribe eth", "err", err)
		return
	}
	logger.Info("finish to subscribe")
	for {
		select {

		case err = <-sub.Err():
			logger.Error("failed to watch eth", "err", err)

		case vLog := <-logs:
			data, err := selectEvent(vLog, oracleCli, functionClient)
			if err != nil {
				continue
			}
			logger.Info("watched log successfully", "data", data)
		}
	}
}

func selectEvent(vLog types.Log, oracleClient *actor.FunctionOracle, functionClient *FunctionConsumer) (interface{}, error) {
	var (
		data interface{}
	)
	logger.Info("log topic", "value", vLog.Topics[0].Hex())
	// TODO: subscribe topic
	switch vLog.Topics[0].Hex() {
	case chain.OracleResponseSignature:
		resp, err := oracleClient.ParseOracleResponse(vLog)
		if err != nil {
			logger.Error("failed to parse function response", "err", err)
			return nil, err
		}
		logger.Info("received OracleResponse event", "log", resp)
		data = resp
	case chain.OracleRequestTimeoutSignature:

		resp, err := oracleClient.ParseOracleRequestTimeout(vLog)
		if err != nil {
			logger.Error("failed to parse request timeout event", "err", err)
			return nil, err
		}
		logger.Info("received OracleRequestTimeout event", "log", resp)
		data = resp
	case chain.RequestFulfilledSignature:
		resp, err := functionClient.ParseRequestFulfilled(vLog)
		if err != nil {
			logger.Error("failed to parse function response", "err", err)
			return nil, err
		}
		logger.Info("received RequestFulfilled event",
			"reqId", hex.EncodeToString(resp.Id[:]),
			"scr node", resp.Node,
			"resp", string(resp.Result),
			"err", string(resp.Err),
		)
		data = resp

	case chain.OracleRequestSignature:

		sent, err := oracleClient.ParseOracleRequest(vLog)
		if err != nil {
			logger.Error("failed to parse OracleRequest event", "err", err)
			return nil, err
		}
		logger.Info("receive OracleRequest event",
			"requestId", hex.EncodeToString(sent.RequestId[:]),
			"requestContract", sent.RequestingContract,
			"requestInitiator", sent.RequestInitiator,
			"from contract", sent.RequestingContract.String(),
			"functionId", hex.EncodeToString(sent.FunctionId[:]),
		)

		logger.Debug("request raw data", "hex req data", hex.EncodeToString(sent.Data))

		reqRawDataMap, err := cbor.ParseDietCBOR(sent.Data)
		if err != nil {
			logger.Error("failed to decode contract request", "err", err)
			return nil, err
		}

		logger.Info("decode requested data", "raw args ", reqRawDataMap)

	default:
		return nil, errors.Errorf("not listen to the event, topic:%s, addr: %v", vLog.Topics[0].Hex(), vLog.Address)
	}
	return data, nil
}
