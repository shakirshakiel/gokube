package registry

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"etcdtest/pkg/api"
	"etcdtest/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

var (
	etcdServer   *embed.Etcd
	etcdClient   *clientv3.Client
	nodeRegistry *NodeRegistry
	ctx          context.Context
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error
	etcdServer, err = startEmbeddedEtcd()
	if err != nil {
		fmt.Printf("Failed to start etcd: %v\n", err)
		os.Exit(1)
	}

	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("Failed to create etcd client: %v\n", err)
		os.Exit(1)
	}

	etcdStorage := storage.NewEtcdStorage(etcdClient)
	nodeRegistry = NewNodeRegistry(etcdStorage)
	ctx = context.Background()
}

func teardown() {
	etcdServer.Close()
	etcdClient.Close()
	os.RemoveAll(etcdServer.Config().Dir)
}

func TestCreateNode(t *testing.T) {
	node := createTestNode("test-node-1", "123")
	err := nodeRegistry.CreateNode(ctx, node)
	assert.NoError(t, err)
}

func TestGetNode(t *testing.T) {
	nodeName := "test-node-2"
	createTestNodeInRegistry(t, nodeName, "456")

	node, err := nodeRegistry.GetNode(ctx, nodeName)
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, nodeName, node.Name)
	assert.Equal(t, "456", node.UID)
	assert.False(t, node.Spec.Unschedulable)
}

func TestGetNonExistentNode(t *testing.T) {
	_, err := nodeRegistry.GetNode(ctx, "non-existent-node")
	assert.Error(t, err)
}

func TestUpdateNode(t *testing.T) {
	nodeName := "test-node-3"
	createTestNodeInRegistry(t, nodeName, "789")

	node, err := nodeRegistry.GetNode(ctx, nodeName)
	require.NoError(t, err)

	node.Spec.Unschedulable = true
	err = nodeRegistry.UpdateNode(ctx, node)
	assert.NoError(t, err)

	updatedNode, err := nodeRegistry.GetNode(ctx, nodeName)
	assert.NoError(t, err)
	assert.True(t, updatedNode.Spec.Unschedulable)
}

func TestListNodes(t *testing.T) {
	// Clear existing nodes
	clearNodes(t)

	// Create test nodes
	createTestNodeInRegistry(t, "test-node-4", "101")
	createTestNodeInRegistry(t, "test-node-5", "102")

	nodes, err := nodeRegistry.ListNodes(ctx)
	assert.NoError(t, err)
	assert.Len(t, nodes, 2)
	assert.Contains(t, []string{nodes[0].Name, nodes[1].Name}, "test-node-4")
	assert.Contains(t, []string{nodes[0].Name, nodes[1].Name}, "test-node-5")
}

func TestDeleteNode(t *testing.T) {
	nodeName := "test-node-6"
	createTestNodeInRegistry(t, nodeName, "103")

	err := nodeRegistry.DeleteNode(ctx, nodeName)
	assert.NoError(t, err)

	_, err = nodeRegistry.GetNode(ctx, nodeName)
	assert.Error(t, err)
}

func TestDeleteNonExistentNode(t *testing.T) {
	err := nodeRegistry.DeleteNode(ctx, "non-existent-node")
	assert.NoError(t, err) // Deleting a non-existent node should not return an error
}

// Helper functions

func createTestNode(name, uid string) *api.Node {
	return &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: name,
			UID:  uid,
		},
		Spec: api.NodeSpec{
			Unschedulable: false,
		},
	}
}

func createTestNodeInRegistry(t *testing.T, name, uid string) {
	node := createTestNode(name, uid)
	err := nodeRegistry.CreateNode(ctx, node)
	require.NoError(t, err)
}

func clearNodes(t *testing.T) {
	nodes, err := nodeRegistry.ListNodes(ctx)
	require.NoError(t, err)
	for _, node := range nodes {
		err := nodeRegistry.DeleteNode(ctx, node.Name)
		require.NoError(t, err)
	}
}

func startEmbeddedEtcd() (*embed.Etcd, error) {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	cfg.LogLevel = "error"
	cfg.LogOutputs = []string{"stderr"}

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}

	select {
	case <-e.Server.ReadyNotify():
		fmt.Println("Embedded etcd is ready!")
	case <-time.After(10 * time.Second):
		e.Server.Stop() // trigger a shutdown
		return nil, fmt.Errorf("server took too long to start")
	}

	return e, nil
}
