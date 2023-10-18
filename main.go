// Copyright (c) OpenFaaS Author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/openfaas/of-watchdog/chain"
	"github.com/openfaas/of-watchdog/executor"

	limiter "github.com/openfaas/faas-middleware/concurrency-limiter"
	"github.com/openfaas/of-watchdog/config"
	"github.com/openfaas/of-watchdog/logger"
	"github.com/openfaas/of-watchdog/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var (
	acceptingConnections int32
)

func main() {
	var runHealthcheck bool
	var versionFlag bool

	flag.BoolVar(&versionFlag, "version", false, "Print the version and exit")
	flag.BoolVar(&runHealthcheck,
		"run-healthcheck",
		false,
		"Check for the a lock-file, when using an exec healthcheck. Exit 0 for present, non-zero when not found.")

	flag.Parse()

	printVersion()

	if versionFlag {
		return
	}

	if runHealthcheck {
		if lockFilePresent() {
			os.Exit(0)
		}

		fmt.Fprintf(os.Stderr, "unable to find lock file.\n")
		os.Exit(1)
	}

	atomic.StoreInt32(&acceptingConnections, 0)

	watchdogConfig, err := config.New(os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err.Error())
		os.Exit(1)
	}

	// baseFunctionHandler is the function invoker without any other middlewares.
	// It is used to provide a generic way to implement the readiness checks regardless
	// of the request mode.
	baseFunctionHandler := buildRequestHandler(watchdogConfig, watchdogConfig.PrefixLogs)
	requestHandler := baseFunctionHandler

	var limit limiter.Limiter
	if watchdogConfig.MaxInflight > 0 {
		requestLimiter := limiter.NewConcurrencyLimiter(requestHandler, watchdogConfig.MaxInflight)
		requestHandler = requestLimiter.Handler()
		limit = requestLimiter
	}

	logger.Info("watch dog info", "Watchdog mode", config.WatchdogMode(watchdogConfig.OperationalMode), "fprocess", watchdogConfig.FunctionProcess)

	chainConfig := config.LoadChainConfig()
	var publisher *chain.Interactor
	defer publisher.Clean()

	if !chainConfig.Disable {
		publisher = chain.NewSubscriber(chainConfig)
		chainHandler := executor.NewChainHandler(publisher, chainConfig.VerifierScoreAddr(), watchdogConfig.UpstreamURL, watchdogConfig.ExecTimeout)
		requestHandler = chainHandler.MakeChainHandler(baseFunctionHandler)
	}

	httpMetrics := metrics.NewHttp()
	http.HandleFunc("/", metrics.InstrumentHandler(requestHandler, httpMetrics))
	http.HandleFunc("/_/health", makeHealthHandler())
	http.Handle("/_/ready", &readiness{
		// make sure to pass original handler, before it's been wrapped by
		// the limiter
		functionHandler: baseFunctionHandler,
		endpoint:        watchdogConfig.ReadyEndpoint,
		lockCheck:       lockFilePresent,
		limiter:         limit,
	})

	metricsServer := metrics.MetricsServer{}
	metricsServer.Register(watchdogConfig.MetricsPort)

	cancel := make(chan bool)

	go metricsServer.Serve(cancel)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", watchdogConfig.TCPPort),
		ReadTimeout:    watchdogConfig.HTTPReadTimeout,
		WriteTimeout:   watchdogConfig.HTTPWriteTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	logger.Info("Timeouts",
		"read", watchdogConfig.HTTPReadTimeout,
		"write", watchdogConfig.HTTPWriteTimeout,
		"hard", watchdogConfig.ExecTimeout,
		"health", watchdogConfig.HealthcheckInterval)

	logger.Info("Listening on port", "p", watchdogConfig.TCPPort)

	listenUntilShutdown(s,
		watchdogConfig.HealthcheckInterval,
		watchdogConfig.HTTPWriteTimeout,
		watchdogConfig.SuppressLock,
		&httpMetrics)
}

func listenUntilShutdown(s *http.Server, healthcheckInterval time.Duration, writeTimeout time.Duration, suppressLock bool, httpMetrics *metrics.Http) {

	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig

		logger.Info("SIGTERM: no new connections in ", "healthcheckInterval", healthcheckInterval.String())

		if err := markUnhealthy(); err != nil {
			logger.Info("Unable to mark server as unhealthy: %s\n", err.Error())
		}

		<-time.Tick(healthcheckInterval)

		connections := int64(testutil.ToFloat64(httpMetrics.InFlight))
		logger.Info("No new connections allowed, draining requests", "draining", connections)

		// The maximum time to wait for active connections whilst shutting down is
		// equivalent to the maximum execution time i.e. writeTimeout.
		ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			logger.Info("Error in Shutdown", "err", err)
		}

		connections = int64(testutil.ToFloat64(httpMetrics.InFlight))

		logger.Info("Exiting. Active connections", "connections", connections)

		close(idleConnsClosed)
	}()

	// Run the HTTP server in a separate go-routine.
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			logger.Info("Error ListenAndServe", "err", err)
			close(idleConnsClosed)
		}
	}()

	if suppressLock == false {
		path, writeErr := createLockFile()

		if writeErr != nil {
			logger.Fatal("cannot write path. To disable lock-file set env suppress_lock=true", "path", path, "err", writeErr.Error())
		}
	} else {
		logger.Info("Warning: \"suppress_lock\" is enabled. No automated health-checks will be in place for your function.")

		atomic.StoreInt32(&acceptingConnections, 1)
	}

	<-idleConnsClosed
}

func getEnvironment(r *http.Request) []string {
	var envs []string

	envs = os.Environ()
	for k, v := range r.Header {
		kv := fmt.Sprintf("Http_%s=%s", strings.Replace(k, "-", "_", -1), v[0])
		envs = append(envs, kv)
	}
	envs = append(envs, fmt.Sprintf("Http_Method=%s", r.Method))

	if len(r.URL.RawQuery) > 0 {
		envs = append(envs, fmt.Sprintf("Http_Query=%s", r.URL.RawQuery))
	}

	if len(r.URL.Path) > 0 {
		envs = append(envs, fmt.Sprintf("Http_Path=%s", r.URL.Path))
	}

	if len(r.TransferEncoding) > 0 {
		envs = append(envs, fmt.Sprintf("Http_Transfer_Encoding=%s", r.TransferEncoding[0]))
	}

	return envs
}

func printVersion() {
	sha := "unknown"
	if len(GitCommit) > 0 {
		sha = GitCommit
	}

	logger.Info("version info", "Version", BuildVersion(), "SHA", sha)
}
