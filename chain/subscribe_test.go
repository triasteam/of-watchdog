package chain

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/openfaas/of-watchdog/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/openfaas/of-watchdog/chain/actor"
	"github.com/openfaas/of-watchdog/logger"
)

func TestParseLog(t *testing.T) {
	t.Skip()

	cfg := config.Chain{
		Id:                 12345678,
		Addr:               "ws://127.0.0.1:9546",
		FunctionClientAddr: "0xe98a2cBE781B4275aFd985E895E92Aea48B235C7",
		FunctionOracleAddr: "0x4B9f0303352a80550455b8323bc9A3D9690ccbDF",
		KeyFilePath:        "",
		KeyPassword:        "123456",
	}
	// key file path
	sub := NewSubscriber(cfg)
	defer sub.Clean()
	sub.ConnectLoop()
	go sub.watch()

	time.Sleep(time.Second * 6000)
}

func TestSubscriber_Send(t *testing.T) {
	t.Skip()

	cfg := config.Chain{
		Id:                 12345678,
		Addr:               "ws://127.0.0.1:9546",
		FunctionClientAddr: "0xe98a2cBE781B4275aFd985E895E92Aea48B235C7",
		FunctionOracleAddr: "0x4B9f0303352a80550455b8323bc9A3D9690ccbDF",
		KeyFilePath:        "",
		KeyPassword:        "123456",
	}
	// key file path
	sub := NewSubscriber(cfg)
	sink := make(chan *actor.FunctionClientRequestFulfilled)
	defer close(sink)
	time.Sleep(time.Second * 2)
	sent, err := sub.functionClient.WatchRequestFulfilled(&bind.WatchOpts{Context: context.Background()}, sink, nil)
	if err != nil {
		return
	}
	defer sent.Unsubscribe()
	data := &actor.FunctionClientRequestFulfilled{}
	select {
	case data = <-sink:
	case err = <-sent.Err():
		t.Error(err.Error())
		return
	}
	logger.Info("receive tx ", "tx hash", data.Raw.TxHash.String(), "reqId", hex.EncodeToString(data.Id[:]))
}
