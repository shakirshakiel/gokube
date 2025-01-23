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

- Basic understanding of Go programming language
- Familiarity with container concepts

## Setup
### 1. Homebrew
Install homebrew by following the instructions from [homebrew website](https://brew.sh/)
### 2. Colima
This project recommends [colima](https://github.com/abiosoft/colima). Feel free to use alterantives like [Racher desktop](https://rancherdesktop.io/), [Docker desktop](https://www.docker.com/products/docker-desktop/),[Podman desktop](https://podman-desktop.io/),[Orbstack](https://orbstack.dev/) etc., if you are already familiar with it.

To install colima, run the following command from the project directory
```bash
make colima/install
```

Once colima is installed, run the following command to start colima immediatley and restart at login:
```bash
brew services start colima
```

Run the following command to start a VM:
```bash
make colima/start
```

Verify colima & docker is working by running the following command:
```bash
docker ps
```

### 3. Go
Install the latest version of golang from the official [website](https://go.dev/doc/install)

Verify golang is installed by running the following command:
```bash
go version
```

### Running Commands

Run the following command to understand what make targets can be run:
```bash
make help
```
### Basic Commands

- To install dependencies - `make deps`
- To format code - `make fmt`
- To run vet - `make vet`
- To run lint - `make lint`
- To run all tests - `make test`
- To run package specific tests(api, controller, kubelet etc.,) - Eg: `make test/api`, `make test/controller`, `make test/kubelet`
- To generate mocks - `make mockgen`
- To build binaries - `make build`
- To build specific binaries - `make build/apiserver`, `make build/controller`, `make build/kubelet`
- To install binaries to GOPATH - `make install`
- To install specific binaries to GOPATH - `make install/apiserver`, `make install/controller`, `make install/kubelet`
- To run all necessary tasks before committing - `make precommit`
- To clean the workspace - `make clean`

## Project Structure

The GoKube project is organized into several key directories:

```
gokube/
├── cmd/
│   ├── apiserver/
│   ├── controller/
│   ├── kubelet/
├── pkg/
│   ├── api/
│   ├── controller/
│   ├── kubelet/
│   ├── listwatch/
│   ├── registry/
│   ├── runtime/
│   ├── scheduler/
│   └── storage/
├── test/
│   └── ...
├── go.mod
├── go.sum
└── README.md

- `pkg/`: Contains the core packages used throughout the project.
  - `api/`: Defines the API objects and clients.
  - `controller/`: Implements the controllers for managing the system state.
  - `kubelet/`: Implements the kubelet functionality.
  - `listwatch/`: Implements the list and watch functionality.
  - `registry/`: Maintains the registry for k8s objects (nodes, pod, replicaset)
  - `runtime/`: Basic runtime utilities
  - `scheduler/`: Implements the scheduling of pods onto nodes.
  - `storage/`: Implements the storage handling via etcd

- `test/`: Contains integration and end-to-end tests.

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
```

# WORK-IN-PROGRESS

## 1. Setting Up the Development Environment

To set up the development environment, there are two alternate options: using `Devbox` or using `limactl`

### Option 1: Using Devbox

1. Install Devbox by following the instructions on the [Devbox GitHub page](https://github.com/jetify-com/devbox).
2. Once Devbox is installed, navigate to the root directory of this project and run:

  ```bash
  devbox shell
  ```

This will automatically install the required packages (`goreleaser` and `lima`) and set up the environment.

### Option 2: Using limactl

If you prefer limactl use the following instructions

1. Install `limactl`:

  ```bash
  brew install lima
  ```

After installing these tools, you can proceed with the rest of the setup instructions.

## Managing the VM

This setup uses the `workbench/debian-12.yaml` configuration and assumes you are running it on an M series MacBook. If you are using a non-M series MacBook, please ask the instructor to provide the necessary instructions.

When the VM is started, it will have all the necessary tools installed, including Docker and etcd. Additionally, the path to the GoKube binary is set, allowing you to run the apiserver, controller, and kubelet directly from the VM shell.

The Makefile includes commands to manage a Lima VM for running GoKube. Here are the instructions to start, stop, delete, and access the VM shell.

### Starting the VMs

To start the VMs, run the following command:

```bash
make start/master
make start/worker1
```

This command will start a Lima instance named `gokube` using the configuration specified in `workbench/debian-12.yaml`.

### Stopping the VM

To stop the VMs, run:

```bash
make stop/master
make stop/worker1
```

This command will stop the `gokube` Lima instance.

### Deleting the VM

To delete the VM, use:

```bash
make delete/master
make delete/worker1
```

This command will delete the `gokube` Lima instance.

### Accessing the VM Shell

To access the shell of the running VM, execute:

```bash
make shell/master
make shell/worker1
```

This command will open a shell in the `gokube` Lima instance, allowing you to interact with the VM directly.