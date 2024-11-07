package registry

import (
	"context"
	"errors"
	"fmt"
	"path"

	"gokube/pkg/api"
	"gokube/pkg/storage"
)

const (
	nodePrefix = "/registry/nodes/"
)

var (
	ErrNodeNotFound      = errors.New("node not found")
	ErrNodeAlreadyExists = errors.New("node already exists")
	ErrListNodesFailed   = errors.New("failed to list nodes")
)

// NodeRegistry provides CRUD operations for Node objects
type NodeRegistry struct {
	storage storage.Storage
}

// NewNodeRegistry creates a new NodeRegistry
func NewNodeRegistry(storage storage.Storage) *NodeRegistry {
	return &NodeRegistry{storage: storage}
}

// generateKey generates the storage key for a given node name
func generateKey(prefix, name string) string {
	return path.Join(prefix, name)
}

// CreateNode stores a new Node
func (r *NodeRegistry) CreateNode(ctx context.Context, node *api.Node) error {
	key := generateKey(nodePrefix, node.Name)
	existingNode := &api.Node{}

	err := r.storage.Get(ctx, key, existingNode)
	if err == nil {
		return fmt.Errorf("%w: %s", ErrNodeAlreadyExists, node.Name)
	}

	return r.storage.Create(ctx, key, node)
}

// GetNode retrieves a Node by name
func (r *NodeRegistry) GetNode(ctx context.Context, name string) (*api.Node, error) {
	key := generateKey(nodePrefix, name)
	node := &api.Node{}

	if err := r.storage.Get(ctx, key, node); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNodeNotFound, name)
	}

	return node, nil
}

// UpdateNode updates an existing Node
func (r *NodeRegistry) UpdateNode(ctx context.Context, node *api.Node) error {
	key := generateKey(nodePrefix, node.Name)
	return r.storage.Update(ctx, key, node)
}

// DeleteNode removes a Node by name
func (r *NodeRegistry) DeleteNode(ctx context.Context, name string) error {
	key := generateKey(nodePrefix, name)
	return r.storage.Delete(ctx, key)
}

// ListNodes retrieves all Nodes
func (r *NodeRegistry) ListNodes(ctx context.Context) ([]*api.Node, error) {
	nodes := make([]*api.Node, 0)

	if err := r.storage.List(ctx, nodePrefix, &nodes); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrListNodesFailed, err)
	}

	return nodes, nil
}
