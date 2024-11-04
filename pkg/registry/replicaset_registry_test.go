package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"

	"gokube/pkg/api"
	"gokube/pkg/storage"
)

func createTestReplicaSet(name string, replicas int32, image string) *api.ReplicaSet {
	return &api.ReplicaSet{
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Spec: api.ReplicaSetSpec{
			Replicas: replicas,
			Selector: map[string]string{"app": "test"},
			Template: api.PodTemplateSpec{
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Image: image,
						},
					},
				},
			},
		},
	}
}

func TestReplicaSetRegistry_Create(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		ctx := context.Background()
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		rs := createTestReplicaSet("test-replicaset", 3, "nginx:latest")
		registry := NewReplicaSetRegistry(etcdStorage)

		err := registry.Create(ctx, rs)
		require.NoError(t, err, "Failed to create ReplicaSet")

		_, err = registry.Get(ctx, "test-replicaset")
		require.NoError(t, err, "Failed to get created ReplicaSet")
	})
}

func TestReplicaSetRegistry_Get(t *testing.T) {
	t.Run("should return ReplicaSet if it exists", func(t *testing.T) {
		storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
			ctx := context.Background()
			etcdStorage := storage.NewEtcdStorage(etcdServer)
			rs := createTestReplicaSet("test-replicaset", 3, "nginx:latest")
			registry := NewReplicaSetRegistry(etcdStorage)

			err := registry.Create(ctx, rs)
			require.NoError(t, err, "Failed to create ReplicaSet")

			retrievedRS, err := registry.Get(ctx, "test-replicaset")
			require.NoError(t, err, "Failed to get ReplicaSet")

			assert.Equal(t, "test-replicaset", retrievedRS.Name)
			assert.Equal(t, int32(3), retrievedRS.Spec.Replicas)
			assert.Len(t, retrievedRS.Spec.Template.Spec.Containers, 1)
			assert.Equal(t, "nginx:latest", retrievedRS.Spec.Template.Spec.Containers[0].Image)
		})
	})

	t.Run("should return error if ReplicaSet does not exist", func(t *testing.T) {
		storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
			etcdStorage := storage.NewEtcdStorage(etcdServer)
			registry := NewReplicaSetRegistry(etcdStorage)
			ctx := context.Background()

			_, err := registry.Get(ctx, "non-existent-replicaset")
			assert.Error(t, err, "Expected error when getting non-existent ReplicaSet")
		})
	})
}

func TestReplicaSetRegistry_Update(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)

		ctx := context.Background()
		registry := NewReplicaSetRegistry(etcdStorage)
		rs := createTestReplicaSet("test-replicaset", 3, "nginx:latest")
		require.NoError(t, registry.Create(ctx, rs))

		updatedRS := createTestReplicaSet("test-replicaset", 5, "nginx:1.19")
		err := registry.Update(ctx, updatedRS)
		require.NoError(t, err, "Failed to update ReplicaSet")

		retrievedRS, err := registry.Get(ctx, "test-replicaset")
		require.NoError(t, err, "Failed to get updated ReplicaSet")

		assert.Equal(t, int32(5), retrievedRS.Spec.Replicas)
		assert.Len(t, retrievedRS.Spec.Template.Spec.Containers, 1)
		assert.Equal(t, "nginx:1.19", retrievedRS.Spec.Template.Spec.Containers[0].Image)
	})
}

func TestReplicaSetRegistry_List(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewReplicaSetRegistry(etcdStorage)
		ctx := context.Background()

		replicaSets := []*api.ReplicaSet{
			createTestReplicaSet("test-replicaset-1", 3, "nginx:latest"),
			createTestReplicaSet("test-replicaset-2", 2, "nginx:1.19"),
		}

		for _, rs := range replicaSets {
			err := registry.Create(ctx, rs)
			require.NoError(t, err)
		}

		rsList, err := registry.List(ctx)
		require.NoError(t, err, "Failed to list ReplicaSets")

		assert.Len(t, rsList, len(replicaSets))
		assert.ElementsMatch(t, replicaSets, rsList)
	})
}

func TestReplicaSetRegistry_Delete(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewReplicaSetRegistry(etcdStorage)
		ctx := context.Background()

		rs := createTestReplicaSet("test-replicaset", 3, "nginx:latest")
		require.NoError(t, registry.Create(ctx, rs))

		err := registry.Delete(ctx, "test-replicaset")
		require.NoError(t, err, "Failed to delete ReplicaSet")

		_, err = registry.Get(ctx, "test-replicaset")
		assert.Error(t, err, "Expected error when getting deleted ReplicaSet")
	})
}
