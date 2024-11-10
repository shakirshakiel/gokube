package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gokube/pkg/kubelet"

	"github.com/spf13/cobra"
)

var (
	nodeName     string
	apiServerURL string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kubelet",
		Short: "Start the gokube kubelet",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runKubelet(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVar(&nodeName, "node-name", "", "Name of this kubelet node")
	rootCmd.Flags().StringVar(&apiServerURL, "api-server", "localhost:8080", "URL of the API server")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runKubelet() error {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	k, err := kubelet.NewKubelet(nodeName, apiServerURL)
	if err != nil {
		return fmt.Errorf("failed to create kubelet: %v", err)
	}

	if err := k.Start(); err != nil {
		return fmt.Errorf("failed to start kubelet: %v", err)
	}

	<-stopCh
	fmt.Println("\nReceived shutdown signal. Stopping kubelet...")
	return nil
}
