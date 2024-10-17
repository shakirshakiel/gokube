package registry

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"etcdtest/pkg/api"
	"etcdtest/pkg/storage"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func setupEtcdStorage() storage.Storage {

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Printf("Failed to create etcd client: %v\n", err)
		os.Exit(1)
	}

	etcdStorage := storage.NewEtcdStorage(etcdClient)
	return etcdStorage
}

func TestPodRegistry_CreateAndUpdatePod(t *testing.T) {
	etcdStorage := setupEtcdStorage()
	registry := NewPodRegistry(etcdStorage)

	ctx := context.Background()

	// Test Create
	pod := &api.Pod{
		Name:   "test-pod",
		Status: api.PodStatusUnassigned,
	}

	err := registry.CreatePod(ctx, pod)
	if err != nil {
		t.Fatalf("Failed to create pod: %v", err)
	}

	// Verify pod was created
	createdPod, err := registry.GetPod(ctx, "test-pod")
	if err != nil {
		t.Fatalf("Failed to get created pod: %v", err)
	}
	if createdPod.Name != "test-pod" || createdPod.Status != api.PodStatusUnassigned {
		t.Errorf("Created pod does not match expected: got %v, want %v", createdPod, pod)
	}

	// Test Update
	updatedPod := &api.Pod{
		Name:   "test-pod",
		Status: api.PodStatusAssigned,
	}

	err = registry.UpdatePod(ctx, updatedPod)
	if err != nil {
		t.Fatalf("Failed to update pod: %v", err)
	}

	// Verify pod was updated
	retrievedPod, err := registry.GetPod(ctx, "test-pod")
	if err != nil {
		t.Fatalf("Failed to get updated pod: %v", err)
	}
	if retrievedPod.Status != api.PodStatusAssigned {
		t.Errorf("Updated pod status does not match expected: got %v, want %v", retrievedPod.Status, api.PodStatusAssigned)
	}

	// Clean up
	err = registry.DeletePod(ctx, "test-pod")
	if err != nil {
		t.Fatalf("Failed to delete pod: %v", err)
	}
}
