package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"

	"gokube/pkg/api"
	"gokube/pkg/storage"
)

func TestNewNodeRegistry(t *testing.T) {
	etcdStorage := storage.NewEtcdStorage(nil)
	nodeRegistry := NewNodeRegistry(etcdStorage)

	assert.NotNil(t, nodeRegistry)
	assert.Equal(t, etcdStorage, nodeRegistry.storage)
}

func TestNodeRegistry_CreateNode(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		nodeRegistry := NewNodeRegistry(etcdStorage)
		node := createTestNode("test-node-1", "123")

		err := nodeRegistry.CreateNode(context.Background(), node)

		assert.NoError(t, err)
	})
}

func TestNodeRegistry_GetNode(t *testing.T) {
	t.Run("should return node if it exists", func(t *testing.T) {
		storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
			etcdStorage := storage.NewEtcdStorage(etcdServer)
			nodeName := "test-node-2"
			nodeRegistry := NewNodeRegistry(etcdStorage)
			ctx := context.Background()

			createTestNodeInRegistry(t, nodeRegistry, nodeName, "456")

			node, err := nodeRegistry.GetNode(ctx, nodeName)
			assert.NoError(t, err)
			assert.NotNil(t, node)
			assert.Equal(t, nodeName, node.Name)
			assert.Equal(t, "456", node.UID)
			assert.False(t, node.Spec.Unschedulable)
		})
	})

	t.Run("should return error if node does not exist", func(t *testing.T) {
		storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
			etcdStorage := storage.NewEtcdStorage(etcdServer)
			nodeRegistry := NewNodeRegistry(etcdStorage)
			ctx := context.Background()

			_, err := nodeRegistry.GetNode(ctx, "non-existent-node")
			assert.Errorf(t, err, "node non-existent-node not found")
		})
	})
}

func TestNodeRegistry_UpdateNode(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		nodeRegistry := NewNodeRegistry(etcdStorage)
		nodeName := "test-node-3"
		createTestNodeInRegistry(t, nodeRegistry, nodeName, "789")

		node, err := nodeRegistry.GetNode(context.Background(), nodeName)
		require.NoError(t, err)

		node.Spec.Unschedulable = true
		err = nodeRegistry.UpdateNode(context.Background(), node)
		assert.NoError(t, err)

		updatedNode, err := nodeRegistry.GetNode(context.Background(), nodeName)
		assert.NoError(t, err)
		assert.True(t, updatedNode.Spec.Unschedulable)
	})
}

func TestNodeRegistry_ListNodes(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		nodeRegistry := NewNodeRegistry(etcdStorage)
		ctx := context.Background()

		// Clear existing nodes
		clearNodes(t, nodeRegistry)

		// Create test nodes
		createTestNodeInRegistry(t, nodeRegistry, "test-node-4", "101")
		createTestNodeInRegistry(t, nodeRegistry, "test-node-5", "102")

		nodes, err := nodeRegistry.ListNodes(ctx)
		assert.NoError(t, err)
		assert.Len(t, nodes, 2)
		assert.Contains(t, []string{nodes[0].Name, nodes[1].Name}, "test-node-4")
		assert.Contains(t, []string{nodes[0].Name, nodes[1].Name}, "test-node-5")
	})
}

func TestNodeRegistry_DeleteNode(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		nodeRegistry := NewNodeRegistry(etcdStorage)
		ctx := context.Background()

		nodeName := "test-node-6"
		createTestNodeInRegistry(t, nodeRegistry, nodeName, "103")

		err := nodeRegistry.DeleteNode(ctx, nodeName)
		assert.NoError(t, err)

		_, err = nodeRegistry.GetNode(ctx, nodeName)
		assert.Error(t, err)
	})
}

func TestDeleteNonExistentNode(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		nodeRegistry := NewNodeRegistry(etcdStorage)
		ctx := context.Background()

		err := nodeRegistry.DeleteNode(ctx, "non-existent-node")
		assert.NoError(t, err) // Deleting a non-existent node should not return an error
	})
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

func createTestNodeInRegistry(t *testing.T, nodeRegistry *NodeRegistry, name, uid string) {
	node := createTestNode(name, uid)
	err := nodeRegistry.CreateNode(context.Background(), node)

	assert.NoError(t, err)

	require.NoError(t, err)
}

func clearNodes(t *testing.T, nodeRegistry *NodeRegistry) {
	ctx := context.Background()
	nodes, err := nodeRegistry.ListNodes(ctx)
	require.NoError(t, err)
	for _, node := range nodes {
		err := nodeRegistry.DeleteNode(ctx, node.Name)
		require.NoError(t, err)
	}
}
