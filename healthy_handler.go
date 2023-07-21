package main

import (
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/openfaas/of-watchdog/logger"
)

func makeHealthHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if atomic.LoadInt32(&acceptingConnections) == 0 || lockFilePresent() == false {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func markUnhealthy() error {
	atomic.StoreInt32(&acceptingConnections, 0)

	path := filepath.Join(os.TempDir(), ".lock")
	logger.Info("Removing lock-file", "path", path)
	removeErr := os.Remove(path)
	return removeErr
}
