// Basic Functionality: List
package listwatch

import (
	"context"
)

// ListFunc defines a function that performs listing
type ListFunc func(ctx context.Context) ([]interface{}, error)

// ListWatch combines list and watch functionality
type ListWatch struct {
	listFn      ListFunc
	watchPrefix string
}

// NewListWatch creates a new ListWatch
func NewListWatch(listFn ListFunc, watchPrefix string) *ListWatch {
	return &ListWatch{
		listFn:      listFn,
		watchPrefix: watchPrefix,
	}
}

// List performs the list operation
func (lw *ListWatch) List(ctx context.Context) ([]interface{}, error) {
	return lw.listFn(ctx)
}

// WatchPrefix returns the prefix to watch
func (lw *ListWatch) WatchPrefix() string {
	return lw.watchPrefix
}
