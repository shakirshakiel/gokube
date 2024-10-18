package kubelet

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestStartContainerWithRealDocker(t *testing.T) {
	// Skip this test if we're not in an environment where we can connect to Docker
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Skip("Skipping test: unable to connect to Docker")
	}

	ctx := context.Background()
	containerName := "test-container"
	imageName := "nginx"

	// Ensure the container doesn't exist before we start
	_ = dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})

	kubelet, err := NewKubelet("test-node", "http://fake-api-server-url")

	if err != nil {
		t.Fatalf("Failed to create Kubelet: %v", err)
	}

	err = kubelet.StartContainer(ctx, containerName, imageName)
	if err != nil {
		t.Fatalf("StartContainer failed: %v", err)
	}

	// Check if the container is running
	containerJSON, err := dockerClient.ContainerInspect(ctx, containerName)
	if err != nil {
		t.Fatalf("Failed to inspect container: %v", err)
	}

	if !containerJSON.State.Running {
		t.Errorf("Container is not running")
	}

	// Clean up: stop and remove the container
	timeout := 10
	err = dockerClient.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout})
	if err != nil {
		t.Errorf("Failed to stop container: %v", err)
	}

	err = dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	if err != nil {
		t.Errorf("Failed to remove container: %v", err)
	}
}
