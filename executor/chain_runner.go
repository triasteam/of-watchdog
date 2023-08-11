package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/docker/go-units"
	"github.com/pkg/errors"

	"github.com/openfaas/of-watchdog/chain"
	"github.com/openfaas/of-watchdog/logger"
)

type ChainHandler struct {
	publisher            chain.Publish
	FunctionsProviderURL string
	Client               *http.Client
	ExecTimeout          time.Duration
}

func NewChainHandler(
	publisher chain.Publish,
	FunctionsProviderURL string,
	timeout time.Duration,
) *ChainHandler {
	return &ChainHandler{
		publisher:            publisher,
		FunctionsProviderURL: FunctionsProviderURL,
		Client:               makeProxyClient(timeout),
		ExecTimeout:          timeout,
	}
}

func (ch ChainHandler) Run() {
	reqData := &chain.FunctionRequest{}
	for true {
		select {
		case dataByte := <-ch.publisher.Receive():

			logger.Info("receive request data", "value", string(dataByte))
			err := json.Unmarshal(dataByte, reqData)
			if err != nil {
				logger.Error("failed to unmarshal from the chain", "data", string(dataByte), "err", err)
				continue
			}
		}
		var errRet []byte
		marshal, err := json.Marshal(reqData.Body)
		if err != nil {
			logger.Error("failed to unmarshal function request body", "err", err)
			continue
		}

		reqHttp, err := http.NewRequest(http.MethodPost, ch.FunctionsProviderURL, bytes.NewReader(marshal))
		if err != nil {
			logger.Error("failed to new http request", "err", err)
			continue
		}
		reqHttp.Header.Set("requestId", reqData.ReqId)
		resp, err := ch.ExecFunction(reqHttp)
		if err != nil {
			errRet = []byte(err.Error())
		}

		ret := &chain.FulFilledRequest{
			RequestId: reqData.ReqId,
			Resp:      resp,
			Err:       errRet,
		}
		ch.publisher.Reply(ret)
		if ret.Err != nil {
			logger.Error("failed to execute function", "err", string(ret.Err))
			continue
		}

		logger.Info("get result from function", "ret", ret)
	}

}
func (ch ChainHandler) ExecFunction(r *http.Request) ([]byte, error) {
	startedTime := time.Now()
	var reqCtx context.Context
	var cancel context.CancelFunc

	if ch.ExecTimeout.Nanoseconds() > 0 {
		reqCtx, cancel = context.WithTimeout(r.Context(), ch.ExecTimeout)
	} else {
		reqCtx = r.Context()
		cancel = func() {}
	}
	defer cancel()

	logger.Info("post function request", "url", r.URL)

	res, err := ch.Client.Do(r.WithContext(reqCtx))
	if err != nil {
		logger.Info("Upstream HTTP request error", "err", err.Error())

		// Error unrelated to context / deadline
		if reqCtx.Err() == nil {
			logger.Error("failed to exec function", "err", err, "cost time", time.Since(startedTime).Seconds())
			return nil, err
		}

		<-reqCtx.Done()

		if reqCtx.Err() != nil {
			// Error due to timeout / deadline

			logger.Error("failed to exec function", "err", err, "timeout", ch.ExecTimeout)
			return nil, errors.WithMessagef(err, "due to exec_timeout,%v", ch.ExecTimeout)
		}

		return nil, err
	}

	done := time.Since(startedTime)

	var bodyBytes []byte
	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			log.Println("read body err", err)
		}
	}
	logger.Info("from chain event", r.Method, "RequestURI", r.RequestURI, "Status", res.Status, "ContentLength", units.HumanSize(float64(res.ContentLength)), "cost time", done.Seconds())
	return bodyBytes, nil
}
