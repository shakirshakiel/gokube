# GoKube: A Miniature Kubernetes-like Container Orchestrator

GoKube is an educational project that implements a simplified version of a container orchestrator, inspired by Kubernetes. This project is designed to teach the concepts of distributed system design using a Kubernetes-like system as an example.

## Project Overview

GoKube is built in Go and aims to demonstrate key concepts of container orchestration such as:

- Container scheduling
- Service discovery
- Load balancing
- State management
- Scaling

By implementing a miniature version of Kubernetes, this project provides hands-on experience with the fundamental principles of distributed systems and container orchestration.

## Prerequisites

- Docker installed on your system
- Basic understanding of Go programming language
- Familiarity with container concepts

## Building the Docker Image

To build the Docker image for development and testing, run the following command from the root directory of this project:

```bash
docker build --no-cache -t gokube-builder .
```
## Running the Development Environment

To run a container from the built image, use the following command:

```bash
docker run --ulimit nofile=65536:65536 --privileged --rm -it -v $(pwd):/app gokube-builder:latest
```

This command does the following:

- `--ulimit nofile=65536:65536`: Sets the ulimit for open files to 65536.
- `--privileged`: Gives extended privileges to this container, necessary for running Docker inside Docker.
- `--rm`: Automatically removes the container when it exits.
- `-it`: Runs the container interactively and allocates a pseudo-TTY.
- `-v $(pwd):/app`: Mounts the current directory to `/app` in the container.
- `gokube-builder:latest`: Specifies the image to use.

Once the container is running, you can build the project using the following command:

```bash
make ci
```
This command will run the continuous integration build process.

### Running Go Commands

Inside the container, you can run all standard Go commands. For example, to run tests for the kubelet package:
```bash
go test -v ./pkg/kubelet
```

You can use similar commands to run tests for other packages, build the project, or perform any other Go-related tasks.

## Project Structure

The GoKube project is organized into several key directories:

```
gokube/
├── cmd/
│   ├── apiserver/
│   ├── controller-manager/
│   ├── kubelet/
│   └── scheduler/
├── pkg/
│   ├── api/
│   ├── controller/
│   ├── kubelet/
│   ├── scheduler/
│   └── util/
├── internal/
│   └── ...
├── test/
│   └── ...
├── docs/
│   └── ...
├── Dockerfile
├── go.mod
├── go.sum
└── README.md

- `pkg/`: Contains the core packages used throughout the project.
  - `api/`: Defines the API objects and clients.
  - `controller/`: Implements the controllers for managing the system state.
  - `kubelet/`: Implements the kubelet functionality.
  - `scheduler/`: Implements the scheduling of pods onto nodes.
  - `util/`: Contains utility functions used across the project.

- `internal/`: Houses internal packages not intended for use outside the project.

- `test/`: Contains integration and end-to-end tests.

- `docs/`: Project documentation.

This structure mimics Kubernetes' organization, providing a familiar layout for those acquainted with the Kubernetes codebase while simplifying it for educational purposes.


## Components

- API Server: Handles API requests and manages the system's state
- Kubelet: Manages containers on individual nodes
- Etcd: Distributed key-value store for system state (simulated)

## Current Features

- Basic container management (create, start, stop)
- Simple pod creation and management
- Rudimentary node management

## TODOs

The following features are planned for implementation:

1. [ ]Implement ReplicationController to create a specified number of replicas of created pods
2. [ ]Implement Scheduler to assign pods to nodes
3. [ ]Implement PodStatus update. Nodes should update pod status periodically with the apiserver
4. [ ]Update the ReplicationController to create newer instances of pods assigned to other nodes if a pod or node hosting the pods fails
5. [ ]Implement a Proxy service to load balance requests across pod instances~


## Learning Objectives

By working with this project, you will gain insights into:

1. The architecture of container orchestration systems
2. Distributed system design principles
3. Container lifecycle management
4. Network management in containerized environments
5. Challenges in distributed state management
6. Scaling and load balancing in distributed systems

## Acknowledgments

- Kubernetes project for inspiration
- Patterns Of Distributed Systems for design principles