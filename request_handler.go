package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/openfaas/of-watchdog/config"
	"github.com/openfaas/of-watchdog/executor"
	"github.com/openfaas/of-watchdog/logger"
)

func buildRequestHandler(watchdogConfig config.WatchdogConfig, prefixLogs bool) http.Handler {
	var requestHandler http.HandlerFunc

	switch watchdogConfig.OperationalMode {
	case config.ModeStreaming:
		requestHandler = makeStreamingRequestHandler(watchdogConfig, prefixLogs, watchdogConfig.LogBufferSize)
	case config.ModeSerializing:
		requestHandler = makeSerializingForkRequestHandler(watchdogConfig, prefixLogs)
	case config.ModeHTTP:
		requestHandler = makeHTTPRequestHandler(watchdogConfig, prefixLogs, watchdogConfig.LogBufferSize)
	case config.ModeStatic:
		requestHandler = makeStaticRequestHandler(watchdogConfig)
	default:
		log.Panicf("unknown watchdog mode: %d", watchdogConfig.OperationalMode)
	}

	return requestHandler
}

func makeSerializingForkRequestHandler(watchdogConfig config.WatchdogConfig, logPrefix bool) func(http.ResponseWriter, *http.Request) {
	functionInvoker := executor.SerializingForkFunctionRunner{
		ExecTimeout:   watchdogConfig.ExecTimeout,
		LogPrefix:     logPrefix,
		LogBufferSize: watchdogConfig.LogBufferSize,
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var environment []string

		if watchdogConfig.InjectCGIHeaders {
			environment = getEnvironment(r)
		}

		commandName, arguments := watchdogConfig.Process()
		req := executor.FunctionRequest{
			Process:       commandName,
			ProcessArgs:   arguments,
			InputReader:   r.Body,
			ContentLength: &r.ContentLength,
			OutputWriter:  w,
			Environment:   environment,
			RequestURI:    r.RequestURI,
			Method:        r.Method,
			UserAgent:     r.UserAgent(),
		}

		w.Header().Set("Content-Type", watchdogConfig.ContentType)
		err := functionInvoker.Run(req, w)
		if err != nil {
			logger.Error("exception", "err", err)
		}
	}
}

type WriterCounter struct {
	w     io.Writer
	bytes int64
}

func (nc *WriterCounter) setWriter(w io.Writer) {
	nc.w = w
}

func (nc *WriterCounter) Bytes() int64 {
	return nc.bytes
}

func (nc *WriterCounter) Write(p []byte) (int, error) {
	n, err := nc.w.Write(p)
	if err != nil {
		return n, err
	}

	nc.bytes += int64(n)
	return n, err
}

func makeStreamingRequestHandler(watchdogConfig config.WatchdogConfig, prefixLogs bool, logBufferSize int) func(http.ResponseWriter, *http.Request) {
	functionInvoker := executor.StreamingFunctionRunner{
		ExecTimeout:   watchdogConfig.ExecTimeout,
		LogPrefix:     prefixLogs,
		LogBufferSize: logBufferSize,
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var environment []string

		if watchdogConfig.InjectCGIHeaders {
			environment = getEnvironment(r)
		}

		ww := WriterCounter{}
		ww.setWriter(w)
		start := time.Now()
		commandName, arguments := watchdogConfig.Process()
		req := executor.FunctionRequest{
			Process:      commandName,
			ProcessArgs:  arguments,
			InputReader:  r.Body,
			OutputWriter: &ww,
			Environment:  environment,
			RequestURI:   r.RequestURI,
			Method:       r.Method,
			UserAgent:    r.UserAgent(),
		}

		w.Header().Set("Content-Type", watchdogConfig.ContentType)
		err := functionInvoker.Run(req)
		if err != nil {
			log.Println(err.Error())

			// Cannot write a status code to the client because we
			// already have written a header
			done := time.Since(start)
			if !strings.HasPrefix(req.UserAgent, "kube-probe") {
				logger.Info("%s %s - %d - ContentLength: %s (%.4fs)", req.Method, req.RequestURI, http.StatusInternalServerError, units.HumanSize(float64(ww.Bytes())), done.Seconds())
				return
			}
		}

		done := time.Since(start)
		if !strings.HasPrefix(req.UserAgent, "kube-probe") {
			logger.Info("%s %s - %d - ContentLength: %s (%.4fs)", req.Method, req.RequestURI, http.StatusOK, units.HumanSize(float64(ww.Bytes())), done.Seconds())
		}
	}
}

func makeHTTPRequestHandler(watchdogConfig config.WatchdogConfig, prefixLogs bool, logBufferSize int) func(http.ResponseWriter, *http.Request) {
	commandName, arguments := watchdogConfig.Process()
	functionInvoker := executor.HTTPFunctionRunner{
		ExecTimeout:    watchdogConfig.ExecTimeout,
		Process:        commandName,
		ProcessArgs:    arguments,
		BufferHTTPBody: watchdogConfig.BufferHTTPBody,
		LogPrefix:      prefixLogs,
		LogBufferSize:  logBufferSize,
	}

	if len(watchdogConfig.UpstreamURL) == 0 {
		log.Fatal(`For "mode=http" you must specify a valid URL for "http_upstream_url"`)
	}

	urlValue, err := url.Parse(watchdogConfig.UpstreamURL)
	if err != nil {
		log.Fatalf(`For "mode=http" you must specify a valid URL for "http_upstream_url", error: %s`, err)
	}

	functionInvoker.UpstreamURL = urlValue

	logger.Info("Forking: %s, arguments: %s", commandName, arguments)
	functionInvoker.Start()

	return func(w http.ResponseWriter, r *http.Request) {

		req := executor.FunctionRequest{
			Process:      commandName,
			ProcessArgs:  arguments,
			InputReader:  r.Body,
			OutputWriter: w,
		}

		if r.Body != nil {
			defer r.Body.Close()
		}

		if err := functionInvoker.Run(req, r.ContentLength, r, w); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}
}

func makeStaticRequestHandler(watchdogConfig config.WatchdogConfig) http.HandlerFunc {
	if watchdogConfig.StaticPath == "" {
		logger.Fatal(`For mode=static you must specify the "static_path" to serve`)
	}

	logger.Info("Serving files at: %s", watchdogConfig.StaticPath)
	return http.FileServer(http.Dir(watchdogConfig.StaticPath)).ServeHTTP
}
