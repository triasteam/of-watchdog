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

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

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
		LocalVerifierAddress: LocalVerifierAddress,
		Client:               makeProxyClient(timeout),
		ExecTimeout:          timeout,
	}
}

func (ch *ChainHandler) MakeChainHandler(preHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			bodyBytes []byte
			errBytes  []byte
			err       error
		)

		lrw := NewLoggingResponseWriter(w)
		preHandler.ServeHTTP(lrw, r)
		if lrw.statusCode != http.StatusOK {
			// Copy the body over
			_, err = io.CopyBuffer(w, bytes.NewReader(errBytes), nil)

		} else {
			// Copy the body over
			_, err = io.CopyBuffer(w, bytes.NewReader(bodyBytes), nil)

		}
		if err != nil {
			logger.Error("failed to copy response body", "err", err)
			return
		}

		reqId := r.Header.Get("requestId")
		if reqId == "" {
			logger.Info("not found requestId")
			return
		}
		ret := &chain.FulFilledRequest{
			RequestId: reqId,
			Resp:      bodyBytes,
			NodeScore: ch.GetLocalWorkScore(),
			Err:       errBytes,
		}

		ch.publisher.Reply(ret)
		if ret.Err != nil {
			logger.Error("failed to execute function", "err", string(ret.Err))
			return
		}
	}

}

func (ch *ChainHandler) GetLocalWorkScore() int64 {

	type VerifierResp struct {
		Data struct {
			Score int `json:"score"`
			Seek  int `json:"seek"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	reqHttp, err := http.NewRequest(http.MethodGet, ch.LocalVerifierAddress, nil)
	if err != nil {
		logger.Error("failed to creat get worker score request, so score is equal to 0 or last score", "score", ch.lastWorkScore, "err", err)
		return ch.lastWorkScore
	}

	resp, err := ch.ExecFunction(reqHttp)
	if err != nil {
		logger.Error("failed to send worker score request, so score is equal to 0 or last score", "score", ch.lastWorkScore, "err", err)
		return ch.lastWorkScore
	}
	ret := VerifierResp{}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		logger.Error("failed to unmarshal worker score response, so score is equal to 0 or last score", "score", ch.lastWorkScore, "err", err)
		return ch.lastWorkScore
	}
	logger.Info("success to get node score", "verifier resp", ret)
	ch.lastWorkScore = int64(ret.Data.Score)
	return ch.lastWorkScore
}

func (ch *ChainHandler) ExecFunction(r *http.Request) ([]byte, error) {
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
