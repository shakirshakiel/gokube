package api

import (
	"time"
)

// Node is a simplified representation of a Kubernetes Node
type Node struct {
	ObjectMeta `json:"metadata,omitempty"`
	Spec       NodeSpec   `json:"spec,omitempty"`
	Status     NodeStatus `json:"status,omitempty"`
}

// ObjectMeta is minimal metadata that all persisted resources must have
type ObjectMeta struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace,omitempty"`
	UID               string    `json:"uid,omitempty"`
	ResourceVersion   string    `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`
}

// NodeSpec describes the basic attributes of a node
type NodeSpec struct {
	Unschedulable bool   `json:"unschedulable,omitempty"`
	ProviderID    string `json:"providerID,omitempty"`
}

type NodeStatus string

// Define some constants for NodeConditionType and ConditionStatus
const (
	NodeNotReady       NodeStatus = "NotReady"
	NodeReady          NodeStatus = "Ready"
	NodeMemoryPressure NodeStatus = "MemoryPressure"
	NodeDiskPressure   NodeStatus = "DiskPressure"
)
