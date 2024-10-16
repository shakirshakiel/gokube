package storage

import (
	"fmt"
	"os"

	"go.etcd.io/etcd/server/v3/embed"
)

// StartEmbeddedEtcd starts an embedded etcd server and returns the server instance,
// the data directory path, and any error encountered.
func StartEmbeddedEtcd() (*embed.Etcd, string, error) {
	dataDir, err := os.MkdirTemp("", "etcd-data")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	cfg := embed.NewConfig()
	cfg.Dir = dataDir

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to start etcd: %v", err)
	}

	select {
	case <-e.Server.ReadyNotify():
		fmt.Printf("Embedded etcd server is ready in %s\n", dataDir)
	case <-e.Server.StopNotify():
		return nil, "", fmt.Errorf("embedded etcd server stopped")
	}

	return e, dataDir, nil
}

// StopEmbeddedEtcd stops the embedded etcd server and removes the data directory.
func StopEmbeddedEtcd(e *embed.Etcd, dataDir string) {
	e.Close()
	os.RemoveAll(dataDir)
	fmt.Println("Embedded etcd server stopped and data directory removed")
}
