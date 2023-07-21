// Copyright (c) OpenFaaS Author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/openfaas/of-watchdog/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer provides instrumentation for HTTP calls
type MetricsServer struct {
	s    *http.Server
	port int
}

// Register binds a HTTP server to expose Prometheus metrics
func (m *MetricsServer) Register(metricsPort int) {

	m.port = metricsPort

	readTimeout := time.Millisecond * 500
	writeTimeout := time.Millisecond * 500

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	m.s = &http.Server{
		Addr:           fmt.Sprintf(":%d", metricsPort),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
		Handler:        metricsMux,
	}

}

// Serve http traffic in go routine, non-blocking
func (m *MetricsServer) Serve(cancel chan bool) {
	logger.Info("Metrics listening on port", "p", m.port)

	go func() {
		if err := m.s.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("metrics error ListenAndServe: %v\n", err))
		}
	}()

	go func() {
		select {
		case <-cancel:
			logger.Info("metrics server shutdown\n")

			m.s.Shutdown(context.Background())
		}
	}()
}

// InstrumentHandler returns a handler which records HTTP requests
// as they are made
func InstrumentHandler(next http.Handler, _http Http) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		then := promhttp.InstrumentHandlerCounter(_http.RequestsTotal,
			promhttp.InstrumentHandlerDuration(_http.RequestDurationHistogram, next))

		_http.InFlight.Inc()
		defer _http.InFlight.Dec()

		then(w, r)
	}
}
