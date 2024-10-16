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

// NodeStatus is minimal information about the current status of a node
type NodeStatus struct {
	Capacity    ResourceList    `json:"capacity,omitempty"`
	Allocatable ResourceList    `json:"allocatable,omitempty"`
	Conditions  []NodeCondition `json:"conditions,omitempty"`
}

// NodeCondition contains condition information for a node
type NodeCondition struct {
	Type   NodeConditionType `json:"type"`
	Status ConditionStatus   `json:"status"`
}

// NodeConditionType is a valid value for NodeCondition.Type
type NodeConditionType string

// ConditionStatus is the status of a condition
type ConditionStatus string

// ResourceList is a set of (resource name, quantity) pairs
type ResourceList map[string]string

// Define some constants for NodeConditionType and ConditionStatus
const (
	NodeReady          NodeConditionType = "Ready"
	NodeMemoryPressure NodeConditionType = "MemoryPressure"
	NodeDiskPressure   NodeConditionType = "DiskPressure"

	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)
