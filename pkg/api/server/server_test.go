package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"

	"etcdtest/pkg/api"
	"etcdtest/pkg/storage"
)

func setupTestEnvironment(t *testing.T) (*APIServer, *embed.Etcd, func()) {
	// Start embedded etcd using util.StartEtcdServer
	e, dataDir, err := storage.StartEmbeddedEtcd()
	if err != nil {
		t.Fatalf("Failed to start etcd: %v", err)
	}

	// Create etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{e.Config().AdvertiseClientUrls[0].String()},
	})
	if err != nil {
		t.Fatalf("Failed to create etcd client: %v", err)
	}

	// Create storage, registry, and API server
	store := storage.NewEtcdStorage(client)
	apiServer := NewAPIServer(store)

	cleanup := func() {
		client.Close()
		storage.StopEmbeddedEtcd(e, dataDir)
	}

	return apiServer, e, cleanup
}

func TestCreateNode(t *testing.T) {
	apiServer, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	node := &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: "test-node",
		},
		Status: api.NodeReady,
	}

	body, _ := json.Marshal(node)
	req := httptest.NewRequest("POST", "/api/v1/nodes", bytes.NewReader(body))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	resp := httptest.NewRecorder()

	container := restful.NewContainer()
	apiServer.registerRoutes(container)
	container.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var createdNode api.Node
	err := json.Unmarshal(resp.Body.Bytes(), &createdNode)
	assert.NoError(t, err)
	assert.Equal(t, node.Name, createdNode.Name)
	assert.Equal(t, node.Status, createdNode.Status)
}

func TestUpdateNodeStatus(t *testing.T) {
	apiServer, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First, create a node
	node := &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: "test-node",
		},
		Status: api.NodeReady,
	}

	err := apiServer.nodeRegistry.CreateNode(context.Background(), node)
	assert.NoError(t, err)

	// Now, update the node's status
	updatedNode := &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: "test-node",
		},
		Status: api.NodeNotReady,
	}

	body, _ := json.Marshal(updatedNode)
	req := httptest.NewRequest("PUT", "/api/v1/nodes/test-node", bytes.NewReader(body))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	resp := httptest.NewRecorder()

	container := restful.NewContainer()
	apiServer.registerRoutes(container)
	container.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedNode api.Node
	err = json.Unmarshal(resp.Body.Bytes(), &returnedNode)
	assert.NoError(t, err)
	assert.Equal(t, updatedNode.Name, returnedNode.Name)
	assert.Equal(t, updatedNode.Status, returnedNode.Status)
}

func TestCreatePod(t *testing.T) {
	apiServer, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	pod := &api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name: "test-pod",
		},
		Spec: api.PodSpec{
			Replicas: 1,
			Containers: []api.Container{
				{
					Image: "nginx:latest",
				},
			},
		},
		// Note: We don't set the Status field here, as it should be set by the server
	}

	body, _ := json.Marshal(pod)
	req := httptest.NewRequest("POST", "/api/v1/pods", bytes.NewReader(body))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	resp := httptest.NewRecorder()

	container := restful.NewContainer()
	apiServer.registerRoutes(container)
	container.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var createdPod api.Pod
	err := json.Unmarshal(resp.Body.Bytes(), &createdPod)
	assert.NoError(t, err)
	assert.Equal(t, pod.Name, createdPod.Name)
	assert.Equal(t, pod.Spec.Replicas, createdPod.Spec.Replicas)
	assert.Equal(t, len(pod.Spec.Containers), len(createdPod.Spec.Containers))
	assert.Equal(t, pod.Spec.Containers[0].Image, createdPod.Spec.Containers[0].Image)

	// Check that the status is set to Unassigned
	assert.Equal(t, api.PodStatusUnassigned, createdPod.Status)
}
