package listwatch

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListWatch(t *testing.T) {
	// Mock list function
	mockList := func(ctx context.Context) ([]interface{}, error) {
		return []interface{}{"item1", "item2"}, nil
	}

	// Create ListWatch
	lw := NewListWatch(mockList, "/test/prefix/")

	// Test List
	items, err := lw.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "item1", items[0])
	assert.Equal(t, "item2", items[1])

	// Test WatchPrefix
	assert.Equal(t, "/test/prefix/", lw.WatchPrefix())
}
