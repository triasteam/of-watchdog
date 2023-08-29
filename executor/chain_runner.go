package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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
	LocalVerifierAddress string
	lastWorkScore        int64
	Client               *http.Client
	ExecTimeout          time.Duration
}

func NewChainHandler(
	publisher chain.Publish,
	LocalVerifierAddress string,
	FunctionsProviderURL string,
	timeout time.Duration,
) *ChainHandler {
	return &ChainHandler{
		publisher:            publisher,
		FunctionsProviderURL: FunctionsProviderURL,
		LocalVerifierAddress: LocalVerifierAddress,
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
			NodeScore: ch.GetLocalWorkScore(),
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

func (ch ChainHandler) GetLocalWorkScore() int64 {

	type VerifierResp struct {
		Data struct {
			Score int `json:"score"`
			Seek  int `json:"seek"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	score := ch.lastWorkScore
	reqHttp, err := http.NewRequest(http.MethodGet, ch.LocalVerifierAddress, nil)
	if err != nil {
		logger.Error("failed to creat get worker score request, so score is equal to 0 or last score", "score", score, "err", err)
		return score
	}

	resp, err := ch.ExecFunction(reqHttp)
	if err != nil {
		logger.Error("failed to send worker score request, so score is equal to 0 or last score", "score", score, "err", err)
		return score
	}
	ret := VerifierResp{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		logger.Error("failed to unmarshal worker score response, so score is equal to 0 or last score", "score", score, "err", err)
		return score
	}
	logger.Info("success to get node score", "verifier resp", ret)
	ch.lastWorkScore = int64(ret.Data.Score)
	return score
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
			logger.Info("read body err", "err", err)
		}
	}
	logger.Info("success to send http req", "method", r.Method, "RequestURI", r.RequestURI, "Status", res.Status, "ContentLength", units.HumanSize(float64(res.ContentLength)), "cost time", done.Seconds())
	return bodyBytes, err
}
