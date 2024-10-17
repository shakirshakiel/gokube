package kubelet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"etcdtest/pkg/api"
)

type Kubelet struct {
	nodeName     string
	apiServerURL string
	pods         map[string]*api.Pod
}

func NewKubelet(nodeName, apiServerURL string, port int) *Kubelet {
	k := &Kubelet{
		nodeName:     nodeName,
		apiServerURL: apiServerURL,
		pods:         make(map[string]*api.Pod),
	}

	return k
}

func (k *Kubelet) Start() error {
	// Register the node with the API server
	if err := k.registerNode(); err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}

	// TODO: Implement other Kubelet functionality here

	// Start watching for pod assignments
	go k.watchPods()

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

func (k *Kubelet) watchPods() {
	for {
		pods, err := k.getPodAssignments()
		if err != nil {
			log.Printf("Error getting pod assignments: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, pod := range pods {
			if _, exists := k.pods[pod.Name]; !exists {
				log.Printf("New pod assigned: %s", pod.Name)
				k.pods[pod.Name] = pod
				go k.runPod(pod)
			}
		}

		time.Sleep(10 * time.Second) // Poll every 10 seconds
	}
}

func (k *Kubelet) getPodAssignments() ([]*api.Pod, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/pods?nodeName=%s", k.apiServerURL, k.nodeName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pods []*api.Pod
	if err := json.NewDecoder(resp.Body).Decode(&pods); err != nil {
		return nil, err
	}

	return pods, nil
}

func (k *Kubelet) runPod(pod *api.Pod) {
	// Simulate running a pod
	log.Printf("Running pod: %s", pod.Name)
	// In a real implementation, this would involve setting up containers, etc.
}
