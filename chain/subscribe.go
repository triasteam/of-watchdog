package chain

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/openfaas/of-watchdog/chain/actor"
	"github.com/openfaas/of-watchdog/chain/cbor"
	"github.com/openfaas/of-watchdog/logger"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

type Publish interface {
	Send([]byte) error
	Receive() chan []byte
	Reply(data *FulFilledRequest)
}

type Configure interface {
	ChainID() int64
	ChainAddr() string
	Key() *keystore.Key
	FuncClientAddr() string
	FuncOracleClientAddr() string
	FuncName() [32]byte
}

type Subscriber struct {
	Configure
	ethCli         *ethclient.Client
	functionClient *actor.FunctionClient
	oracleClient   *actor.FunctionOracle

	renewChan   chan struct{}
	dailEthDone chan struct{}
	cleanOnce   sync.Once
	isRenewed   atomic.Bool

	publishChannel chan []byte
	repliedChan    chan *FulFilledRequest
}

func NewSubscriber(cfg Configure) *Subscriber {

	sub := &Subscriber{
		Configure: cfg,

		renewChan:      make(chan struct{}),
		dailEthDone:    make(chan struct{}),
		publishChannel: make(chan []byte, 100),
		repliedChan:    make(chan *FulFilledRequest, 100),
		cleanOnce:      sync.Once{},
	}

	go sub.ConnectLoop()
	go sub.FulfillRequest()
	go sub.watch()

	return sub
}

func (cs *Subscriber) Clean() {
	cs.cleanOnce.Do(func() {
		close(cs.renewChan)
		close(cs.dailEthDone)
		cs.ethCli.Close()
	})

}

func (cs *Subscriber) resetEthCli(cli *ethclient.Client) {

	cs.ethCli = cli
	funcCli, err := actor.NewFunctionClient(common.HexToAddress(cs.FuncClientAddr()), cs.ethCli)
	if err != nil {
		logger.Error("failed to new FunctionClient", "err", err)
	}
	oracleCli, err := actor.NewFunctionOracle(common.HexToAddress(cs.FuncOracleClientAddr()), cs.ethCli)
	if err != nil {
		logger.Error("failed to new FunctionOracle", "err", err)
	}
	cs.oracleClient = oracleCli
	cs.functionClient = funcCli
}

func (cs *Subscriber) retryDailEth(addr string) {
	start := time.Now()
	err := retry.Do(
		func() error {

			ethCli, err := ethclient.Dial(addr)
			if err != nil {
				logger.Error("failed to connect to ", "node addr", cs.ChainAddr())
				return err
			}
			cs.resetEthCli(ethCli)
			return nil
		},
		retry.Attempts(0),
		retry.Delay(time.Second),
		retry.MaxDelay(2*time.Second),
	)
	dur := time.Since(start)
	if err != nil {
		logger.Error("retry exception during watching", "err", err)
	}
	logger.Info("re-dail eth cost time", "dur", dur.String())
}

func (cs *Subscriber) ConnectLoop() {
	for {
		select {
		case _, ok := <-cs.renewChan:
			if !ok {
				logger.Error("resubscribe channel is closed")
				panic("resubscribe channel is closed")
			}
			cs.retryDailEth(cs.ChainAddr())
			cs.dailEthDone <- struct{}{}

		}
	}
}

type FulFilledRequest struct {
	RequestId string
	Resp      []byte
	Err       []byte
}

func (cs *Subscriber) FulfillRequest() {
	timer := time.NewTicker(2 * time.Second)
	defer timer.Stop()
	for {
		var ret *FulFilledRequest
		select {
		case <-timer.C:
			logger.Info("waiting to fulfill request")
		case ret = <-cs.repliedChan:
			reqID, err := hex.DecodeString(ret.RequestId)
			if err != nil {
				logger.Error("failed to decode request id", "requestId", ret.RequestId)
				continue
			}
			auth, err := bind.NewKeyedTransactorWithChainID(cs.Key().PrivateKey, new(big.Int).SetInt64(cs.ChainID()))
			if err != nil {
				logger.Error("failed to new keyed tx", "err", err)
				continue
			}
			if len(reqID) != 32 {
				logger.Error("unexpected requestId", "requestId", string(reqID))
				continue
			}
			requestId := [32]byte(reqID)
			err = retry.Do(
				func() error {
					if !cs.isRenewed.Load() {
						return errors.New("subscriber is resubscribing, isRenewed is false")
					}
					tx, err := cs.functionClient.HandleOracleFulfillment(&bind.TransactOpts{
						From:   auth.From,
						Signer: auth.Signer,
					}, requestId, ret.Resp, ret.Err)
					if err != nil {
						logger.Error("failed to call HandleOracleFulfillment", "err", err)
						return err
					}
					logger.Info("fulfilled request", "tx hash", tx.Hash().String())
					return nil
				},
				retry.Attempts(5),
				retry.Delay(100*time.Millisecond),
				retry.MaxDelay(300*time.Millisecond))

			if err != nil {
				logger.Error("failed to call HandleOracleFulfillment", "err", err)
				continue
			}

		}

	}

}

func (cs *Subscriber) Reply(data *FulFilledRequest) {
	cs.repliedChan <- data
}

func (cs *Subscriber) Send(data []byte) error {
	cs.publishChannel <- data
	return nil
}

func (cs *Subscriber) Receive() chan []byte {
	return cs.publishChannel
}

func (cs *Subscriber) watch() {

	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(cs.FuncClientAddr()),
			common.HexToAddress(cs.FuncOracleClientAddr()),
		},
		Topics: [][]common.Hash{},
	}
	logs := make(chan types.Log)
	var (
		sub ethereum.Subscription
		err error
	)

	timer := time.NewTicker(2 * time.Second)
	defer timer.Stop()
	for {
		for !cs.isRenewed.Load() {
			logger.Info("############# not found chain log event subscriber, resubscribe")

			cs.renewChan <- struct{}{}
			select {
			case <-cs.dailEthDone:
				err = retry.Do(
					func() error {
						logger.Info("start subscribe")
						sub, err = cs.ethCli.SubscribeFilterLogs(context.Background(), query, logs)
						if err != nil {
							logger.Error("failed to finish to subscribe eth", "err", err)
							return err
						}
						logger.Info("finish to subscribe")
						cs.isRenewed.CAS(false, true)
						return nil
					},
					retry.Attempts(5),
					retry.Delay(500*time.Millisecond),
					retry.MaxDelay(time.Second),
				)
				if err != nil {
					logger.Error("retry exception during watching", "err", err)
				} else {
					logger.Info("############# finish to resubscribe")
				}
			}
		}

		select {
		case <-timer.C:
			logger.Info("watching event from the chain")
		case err = <-sub.Err():
			logger.Error("failed to watch eth", "err", err)
			sub.Unsubscribe()
			sub = nil
			cs.isRenewed.CAS(true, false)
		case vLog := <-logs:
			data, err := cs.selectEvent(vLog)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			logger.Info("watched event successfully", "data", data)
		}
	}

}

func (cs *Subscriber) selectEvent(vLog types.Log) (interface{}, error) {
	var (
		data interface{}
	)
	// TODO: subscribe topic
	switch vLog.Topics[0].Hex() {

	case RequestFulfilledSignature:
		resp, err := cs.functionClient.ParseRequestFulfilled(vLog)
		if err != nil {
			return nil, err
		}
		logger.Info("parse function response", "resp", string(resp.Result))
		data = resp
	case RequestSentSignature:
		sent, err := cs.functionClient.ParseRequestSent(vLog)
		if err != nil {
			return nil, err
		}
		logger.Info("parse sent function request", "resp", sent.Id)
		data = sent
	case OracleRequestSignature:

		sent, err := cs.oracleClient.ParseOracleRequest(vLog)
		if err != nil {
			return nil, err
		}
		logger.Info("receive request event",
			"requestId", hex.EncodeToString(sent.RequestId[:]),
			"requestContract", sent.RequestingContract,
			"requestInitiator", sent.RequestInitiator)

		logger.Debug("request raw data", "hex req data", hex.EncodeToString(sent.Data))
		nameByte := cs.FuncName()
		if bytes.Compare(nameByte[:], sent.RequestId[:]) != 0 {
			logger.Info("do not call function, its name is different")
			return nil, nil
		}

		reqRawDataMap, err := cbor.ParseDietCBOR(sent.Data)
		if err != nil {
			logger.Error("failed to decode contract request", "err", err)
			return nil, err
		}

		logger.Info("decode requested data", "raw args ", reqRawDataMap)
		err = callFunction(cs, hex.EncodeToString(sent.RequestId[:]), reqRawDataMap)
		if err != nil {
			logger.Error("failed to call function", "err", err)
			return nil, err
		}

	default:
		return nil, errors.Errorf("not support event, topic:%s", vLog.Topics[0].Hex())
	}
	return data, nil
}

type FunctionRequest struct {
	FunctionName string
	RequestURL   string
	ReqId        string
	Body         map[string]interface{}
}

func callFunction(pub Publish, reqID string, reqRawDataMap map[string]interface{}) error {

	fnBodyMap := map[string]interface{}{}

	if v, ok := reqRawDataMap["args"]; ok && v != nil {

		var sets []interface{}
		switch v.(type) {
		case []interface{}:
			sets = v.([]interface{})
		default:
			logger.Error("request body has wrong format", "args", v, "type", reflect.TypeOf(v))
		}

		if !ok {
			logger.Error("request args format is wrong", "args", v, "type", reflect.TypeOf(v).Name())
			return errors.Errorf("request args format is wrong, aegs: %v", v)
		}
		for i := 0; i+1 < len(sets); i += 2 {
			fnBodyMap[fmt.Sprintf("%v", sets[i])] = sets[i+1]
		}
	}

	fr := FunctionRequest{
		ReqId: reqID,
		Body:  fnBodyMap,
	}
	bodyBytes, err := json.Marshal(fr)
	if err != nil {
		logger.Error("failed to decode request args to map", "err", err)
		return err
	}
	_ = pub.Send(bodyBytes)
	return nil
}
