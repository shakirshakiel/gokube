package storage

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

func TestEmbeddedEtcd(t *testing.T) {
	// Step 1: Start embedded etcd
	etcdServer, dataDir, err := StartEmbeddedEtcd()
	if err != nil {
		t.Fatalf("Failed to start embedded etcd: %v", err)
	}
	defer StopEmbeddedEtcd(etcdServer, dataDir)

	// Step 2: Set up etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create etcd client: %v", err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 3: Put a key-value pair into etcd
	_, err = cli.Put(ctx, "test-key", "test-value")
	if err != nil {
		t.Fatalf("Failed to put key-value: %v", err)
	}

	// Step 4: Get the value from etcd
	resp, err := cli.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if len(resp.Kvs) != 1 || string(resp.Kvs[0].Value) != "test-value" {
		t.Fatalf("Expected 'test-value', got '%s'", string(resp.Kvs[0].Value))
	}

	// Step 5: Verify the value
	fmt.Printf("Key: %s, Value: %s\n", resp.Kvs[0].Key, resp.Kvs[0].Value)

	// Additional test cases
	t.Run("UpdateExistingKey", func(t *testing.T) {
		_, err := cli.Put(ctx, "test-key", "updated-value")
		if err != nil {
			t.Fatalf("Failed to update key: %v", err)
		}

		resp, err := cli.Get(ctx, "test-key")
		if err != nil {
			t.Fatalf("Failed to get updated key: %v", err)
		}

		if string(resp.Kvs[0].Value) != "updated-value" {
			t.Fatalf("Expected 'updated-value', got '%s'", string(resp.Kvs[0].Value))
		}
	})

	t.Run("DeleteKey", func(t *testing.T) {
		_, err := cli.Delete(ctx, "test-key")
		if err != nil {
			t.Fatalf("Failed to delete key: %v", err)
		}

		resp, err := cli.Get(ctx, "test-key")
		if err != nil {
			t.Fatalf("Failed to get deleted key: %v", err)
		}

		if len(resp.Kvs) != 0 {
			t.Fatalf("Expected key to be deleted, but it still exists")
		}
	})

	t.Run("ListKeys", func(t *testing.T) {
		// Add multiple keys
		_, err := cli.Put(ctx, "/prefix/key1", "value1")
		if err != nil {
			t.Fatalf("Failed to put key1: %v", err)
		}
		_, err = cli.Put(ctx, "/prefix/key2", "value2")
		if err != nil {
			t.Fatalf("Failed to put key2: %v", err)
		}

		// List keys with prefix
		resp, err := cli.Get(ctx, "/prefix/", clientv3.WithPrefix())
		if err != nil {
			t.Fatalf("Failed to list keys: %v", err)
		}

		if len(resp.Kvs) != 2 {
			t.Fatalf("Expected 2 keys, got %d", len(resp.Kvs))
		}
	})

	t.Run("WatchKey", func(t *testing.T) {
		watchKey := "/watch-test/key"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Start watching the key in a separate goroutine
		watchChan := cli.Watch(ctx, watchKey)

		go func() {
			time.Sleep(1 * time.Second) // Wait a bit before making changes
			_, err := cli.Put(ctx, watchKey, "initial-value")
			if err != nil {
				t.Errorf("Failed to put initial value: %v", err)
			}
			time.Sleep(1 * time.Second)
			_, err = cli.Put(ctx, watchKey, "updated-value")
			if err != nil {
				t.Errorf("Failed to update value: %v", err)
			}
			time.Sleep(1 * time.Second)
			_, err = cli.Delete(ctx, watchKey)
			if err != nil {
				t.Errorf("Failed to delete key: %v", err)
			}
		}()

		expectedEvents := []struct {
			Type  mvccpb.Event_EventType
			Value string
		}{
			{mvccpb.PUT, "initial-value"},
			{mvccpb.PUT, "updated-value"},
			{mvccpb.DELETE, ""},
		}

		for i, expected := range expectedEvents {
			select {
			case watchResp := <-watchChan:
				fmt.Printf("Watch Response %d:\n", i+1)
				fmt.Printf("  CompactRevision: %d\n", watchResp.CompactRevision)
				fmt.Printf("  Created: %v\n", watchResp.Created)
				fmt.Printf("  Canceled: %v\n", watchResp.Canceled)
				fmt.Printf("  Header: %+v\n", watchResp.Header)

				for j, ev := range watchResp.Events {
					fmt.Printf("  Event %d:\n", j+1)
					fmt.Printf("    Type: %v\n", ev.Type)
					fmt.Printf("    Key: %s\n", string(ev.Kv.Key))
					fmt.Printf("    Value: %s\n", string(ev.Kv.Value))
					fmt.Printf("    Version: %d\n", ev.Kv.Version)
					fmt.Printf("    ModRevision: %d\n", ev.Kv.ModRevision)
					if ev.PrevKv != nil {
						fmt.Printf("    PrevValue: %s\n", string(ev.PrevKv.Value))
						fmt.Printf("    PrevVersion: %d\n", ev.PrevKv.Version)
					}
				}
				fmt.Println()

				if len(watchResp.Events) != 1 {
					t.Fatalf("Expected 1 event, got %d", len(watchResp.Events))
				}
				ev := watchResp.Events[0]
				if ev.Type != expected.Type {
					t.Errorf("Event %d: Expected type %v, got %v", i, expected.Type, ev.Type)
				}
				if string(ev.Kv.Value) != expected.Value {
					t.Errorf("Event %d: Expected value '%s', got '%s'", i, expected.Value, string(ev.Kv.Value))
				}
			case <-ctx.Done():
				t.Fatalf("Watch timed out")
			}
		}
	})
}
