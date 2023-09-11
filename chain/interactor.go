package chain

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/openfaas/of-watchdog/chain/actor"
	"github.com/openfaas/of-watchdog/logger"
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
	VerifierScoreAddr() string
}

type Interactor struct {
	Configure
	ethCli         *ethclient.Client
	functionClient *actor.FunctionClient
	oracleClient   *actor.FunctionOracle

	renewChan   chan struct{} //notify for reconnecting to bsc node
	dailEthDone chan struct{} // write dailEthDone after connecting to bsc node
	cleanOnce   sync.Once
	isRenewed   atomic.Bool // set isRenewed to true when to receive error from bsc subscriber, then notify to reconnect to eth

	publishChannel chan []byte
	repliedChan    chan *FulFilledRequest
}

func NewSubscriber(cfg Configure) *Interactor {

	sub := &Interactor{
		Configure: cfg,

		renewChan:      make(chan struct{}),
		dailEthDone:    make(chan struct{}),
		publishChannel: make(chan []byte, 100),
		repliedChan:    make(chan *FulFilledRequest, 100),
		cleanOnce:      sync.Once{},
	}
	logger.Info("function signature", "value", signaturesMap)
	go sub.ConnectLoop()
	go sub.FulfillRequest()
	go sub.watch()

	return sub
}

func (cs *Interactor) Clean() {
	cs.cleanOnce.Do(func() {
		close(cs.renewChan)
		close(cs.dailEthDone)
		cs.ethCli.Close()
	})

}

func (cs *Interactor) resetEthCli(cli *ethclient.Client) {

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

func (cs *Interactor) retryDailEth(addr string) {
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

func (cs *Interactor) ConnectLoop() {
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
	NodeScore int64
	Err       []byte
}

func (cs *Interactor) FulfillRequest() {
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

					sink := make(chan *actor.FunctionOracleOracleResponse)
					defer close(sink)
					respSub, err := cs.oracleClient.WatchOracleResponse(&bind.WatchOpts{Context: context.Background()}, sink, [][32]byte{requestId})
					if err != nil {
						return err
					}
					defer respSub.Unsubscribe()
					logger.Info("start to fulfill request")
					tx, err := cs.oracleClient.FulfillRequestByNode(&bind.TransactOpts{
						From:   auth.From,
						Signer: auth.Signer,
					}, requestId, common.HexToAddress(cs.Configure.FuncClientAddr()), new(big.Int).SetInt64(ret.NodeScore), ret.Resp, ret.Err)
					if err != nil {
						logger.Error("cannot send FulfillRequestByNode tx", "requestId", reqID, "err", err)
						return errors.WithMessagef(err, "cannot send FulfillRequestByNode tx")
					}
					logger.Info("wait to chain log event")
					select {
					case resp := <-sink:
						logger.Info("node fulfilled request", "tx hash", tx.Hash().String(), "blockNumber", resp.Raw.BlockNumber, "tx", resp.Raw.TxHash)

					case err = <-respSub.Err():
						logger.Error("failed to send resp", "err", err)
						return err
					}
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

func (cs *Interactor) Reply(data *FulFilledRequest) {

	cs.repliedChan <- data
}

func (cs *Interactor) Send(data []byte) error {
	cs.publishChannel <- data
	return nil
}

func (cs *Interactor) Receive() chan []byte {
	return cs.publishChannel
}

func (cs *Interactor) watch() {

	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash(RequestFulfilled),
				crypto.Keccak256Hash(OracleResponse),
			},
		},
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
			err = cs.selectEvent(vLog)
			if err != nil {
				continue
			}
		}
	}

}

func (cs *Interactor) selectEvent(vLog types.Log) error {

	logger.Info("log topic", "value", vLog.Topics[0].Hex())
	switch vLog.Topics[0].Hex() {
	case OracleResponseSignature:
		resp, err := cs.oracleClient.ParseOracleResponse(vLog)
		if err != nil {
			logger.Error("failed to parse function response", "err", err)
			return err
		}
		logger.Info("received OracleResponse event", "log", resp)

	case RequestFulfilledSignature:
		resp, err := cs.functionClient.ParseRequestFulfilled(vLog)
		if err != nil {
			logger.Error("failed to parse function response", "err", err)
			return err
		}
		logger.Info("received RequestFulfilled event",
			"reqId", hex.EncodeToString(resp.Id[:]),
			"scr node", resp.Node,
			"score", resp.Score,
			"resp", string(resp.Result),
			"err", string(resp.Err),
		)
	default:
		logger.Info("not listen to the event", "topic", vLog.Topics[0].Hex(), "contract address", vLog.Address)
	}
	return nil
}
