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

var defaultGasPrice = big.NewInt(200)

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
	timer := time.NewTicker(60 * time.Second)
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
			var nonce uint64
			blockNumber, err := cs.ethCli.BlockNumber(context.Background())
			if err != nil {
				logger.Error("failed to get block number", "err", err)
			} else {
				nonce, err = cs.ethCli.NonceAt(context.Background(), cs.Key().Address, big.NewInt(int64(blockNumber)))
				if err != nil {
					logger.Error("failed to get nonce", "err", err)
				}
			}
			gasPrice, err := cs.ethCli.SuggestGasPrice(context.Background())
			if err != nil {
				logger.Error("failed to get suggest gas price", "err", err)
				gasPrice = defaultGasPrice
			}
			requestId := [32]byte(reqID)
			err = retry.Do(
				func() error {
					if !cs.isRenewed.Load() {
						return errors.New("subscriber is resubscribing, isRenewed is false")
					}
					var (
						retryErr error
					)
					gasPrice = gasPrice.Add(gasPrice, big.NewInt(1))
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
					defer cancel()
					//sink := make(chan *actor.FunctionClientRequestFulfilled)
					//defer close(sink)
					//respSub, retryErr := cs.functionClient.WatchRequestFulfilled(&bind.WatchOpts{Context: ctx}, sink, [][32]byte{requestId}, nil)
					//if retryErr != nil {
					//	return retryErr
					//}
					//defer respSub.Unsubscribe()

					logger.Info("start to fulfill request", "gas price", gasPrice.String(), "nonce", nonce)
					tx, retryErr := cs.functionClient.HandleOracleFulfillment(&bind.TransactOpts{
						From:     auth.From,
						Signer:   auth.Signer,
						GasPrice: gasPrice,
						Nonce:    big.NewInt(int64(nonce)),
					}, requestId, new(big.Int).SetInt64(ret.NodeScore), ret.Resp, ret.Err)
					if retryErr != nil {
						logger.Error("cannot send HandleOracleFulfillment tx", "from", auth.From, "gasPrice", gasPrice.String(), "nonce", nonce, "requestId", reqID, "err", retryErr)
						return errors.WithMessagef(retryErr, "cannot send HandleOracleFulfillment tx,from %s,", auth.From)
					}
					mined, retryErr := bind.WaitMined(ctx, cs.ethCli, tx)
					if retryErr != nil {
						return retryErr
					}
					//logger.Info("wait to chain log event")
					//select {
					//case <-ctx.Done():
					//	retryErr = ctx.Err()
					//	logger.Error("wait to HandleOracleFulfillment finish timeout", "err", ctx.Err())
					//case resp := <-sink:
					//	logger.Info("node fulfilled request", "tx hash", tx.Hash().String(), "blockNumber", resp.Raw.BlockNumber, "tx", resp.Raw.TxHash)
					//
					//case retryErr = <-respSub.Err():
					//	logger.Error("failed to send resp", "err", retryErr)
					//	//return err
					//}
					logger.Info("finish waiting to mint tx", "from", auth.From, "number", mined.BlockNumber.String(), "tx hash", mined.TxHash.String())
					return retryErr
				},
				retry.Attempts(5),
				retry.Delay(10*time.Millisecond),
				retry.MaxDelay(60*time.Second),
			)

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
	//cs.publishChannel <- data
	return nil
}

func (cs *Interactor) Receive() chan []byte {
	return nil
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

	timer := time.NewTicker(60 * time.Second)
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
