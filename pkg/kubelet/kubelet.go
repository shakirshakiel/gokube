package kubelet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"etcdtest/pkg/api"
)

type Kubelet struct {
	nodeName     string
	apiServerURL string
}

func NewKubelet(nodeName, apiServerURL string) *Kubelet {
	return &Kubelet{
		nodeName:     nodeName,
		apiServerURL: apiServerURL,
	}
}

func (k *Kubelet) Start() error {
	// Register the node with the API server
	if err := k.registerNode(); err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}

	// TODO: Implement other Kubelet functionality here

	return nil
}

func (k *Kubelet) registerNode() error {
	node := &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: k.nodeName,
		},
		Status: api.NodeReady,
	}

	jsonData, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node data: %w", err)
	}

	resp, err := http.Post(k.apiServerURL+"/api/v1/nodes", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request to API server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register node, status code: %d", resp.StatusCode)
	}

	return nil
}
