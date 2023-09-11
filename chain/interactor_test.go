package chain

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/openfaas/of-watchdog/chain/actor"

	"github.com/openfaas/of-watchdog/config"

	"github.com/openfaas/of-watchdog/logger"
)

func TestParseLog(t *testing.T) {
	t.Skip()

	cfg := config.Chain{
		Id:                 12345678,
		Addr:               "ws://127.0.0.1:9546",
		FunctionClientAddr: "0x18306eF4d08c69F5768f0af53D4Bc20C25c6Ad75",
		FunctionOracleAddr: "0x0000000000000000000000000000000000002004",
		KeyFilePath:        "./testdata/UTC--2023-06-05T09-50-10.886531000Z--989777e983d4fccba32d857d797fdb75c27571c5",
		KeyPassword:        "123456",
		FunctionName:       "test1",
	}
	// key file path
	sub := NewSubscriber(cfg)
	defer sub.Clean()
	//sub.ConnectLoop()
	//go sub.watch()
	//function signature	{"RequestFulfilledSignature": "0xf2fa06652e54791d449ab43ede930a20d3b71ed330cad4018f47ba6cc15da00a", "RequestSentSignature": "0x91f0d67c2f27abd6cfc317e120d5e80b31e97b9926b65d3887e59402fb20adfb", "OracleRequestSignature": "0x8fe1923fd8e0dc61a5bd16b4ed3ede4f2c6ee0de6f729fab847432965b138aa3"}
	time.Sleep(time.Second * 6000)
}

func TestInteractor_FulfillRequest(t *testing.T) {
	data, err := hex.DecodeString("588e0da0be236e399f38e8fb37c51fdcec76b98818ecdcc814efb210141d933d")
	fmt.Println(data, " ", err)
	requestId := [32]byte(data)

	fmt.Println(hex.EncodeToString(requestId[:]))
	fmt.Println(string(requestId[:]))
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
	sent, err := sub.functionClient.WatchRequestFulfilled(&bind.WatchOpts{Context: context.Background()}, sink, [][32]byte{}, nil)
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
