package main

import (
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/openfaas/of-watchdog/logger"
)

// createLockFile returns a path to a lock file and/or an error
// if the file could not be created.
func createLockFile() (string, error) {
	path := filepath.Join(os.TempDir(), ".lock")
	logger.Info("Writing lock-file to", "path", path)

	if err := os.MkdirAll(os.TempDir(), os.ModePerm); err != nil {
		return path, err
	}

	if err := os.WriteFile(path, []byte{}, 0660); err != nil {
		return path, err
	}

	atomic.StoreInt32(&acceptingConnections, 1)
	return path, nil
}

func lockFilePresent() bool {
	path := filepath.Join(os.TempDir(), ".lock")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
