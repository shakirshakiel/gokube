package kubelet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"etcdtest/pkg/api"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type Kubelet struct {
	nodeName     string
	apiServerURL string
	dockerClient *client.Client
	pods         map[string]*api.Pod
}

func NewKubelet(nodeName, apiServerURL string) (*Kubelet, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &Kubelet{
		nodeName:     nodeName,
		apiServerURL: apiServerURL,
		dockerClient: dockerClient,
		pods:         make(map[string]*api.Pod),
	}, nil
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

func (k *Kubelet) StartContainer(ctx context.Context, name, imageName string) error {

	log.Printf("Pulling image: %s", imageName)

	// Pull the image
	out, err := k.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %v", imageName, err)
	}

	log.Printf("Successfully pulled image: %s", "nginx")

	// Create the container
	resp, err := k.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "nginx",
		// You can add more configuration options here as needed
	}, nil, nil, nil, name)
	if err != nil {
		return fmt.Errorf("failed to create container %s: %v", name, err)
	}

	// Start the container
	if err := k.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s: %v", name, err)
	}

	fmt.Printf("Started container %s with ID %s\n", name, resp.ID)
	return nil
}
