package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gokube/pkg/api"
	"gokube/pkg/storage"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestNewPodRegistry(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewPodRegistry(etcdStorage)

		assert.NotNil(t, registry)
		assert.Equal(t, etcdStorage, registry.storage)
	})
}

func TestPodRegistry_GetPod(t *testing.T) {
	t.Run("should return pod if it exists", func(t *testing.T) {
		storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
			etcdStorage := storage.NewEtcdStorage(etcdServer)
			registry := NewPodRegistry(etcdStorage)
			ctx := context.Background()

			pod := &api.Pod{
				ObjectMeta: api.ObjectMeta{
					Name: "test-pod",
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Image: "nginx:latest",
						},
					},
					Replicas: 3,
				},
				Status: api.PodPending,
			}

			err := registry.CreatePod(ctx, pod)
			require.NoError(t, err)

			// Test GetPod
			retrievedPod, err := registry.GetPod(ctx, "test-pod")
			require.NoError(t, err)

			// Verify pod name and status
			assert.Equal(t, "test-pod", retrievedPod.Name)
			assert.Equal(t, api.PodPending, retrievedPod.Status)

			// Verify pod spec
			assert.Len(t, retrievedPod.Spec.Containers, 1)
			assert.Equal(t, "nginx:latest", retrievedPod.Spec.Containers[0].Image)
			assert.Equal(t, int32(3), retrievedPod.Spec.Replicas)
		})
	})

	t.Run("should return error if pod does not exist", func(t *testing.T) {
		storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
			etcdStorage := storage.NewEtcdStorage(etcdServer)
			registry := NewPodRegistry(etcdStorage)
			ctx := context.Background()

			_, err := registry.GetPod(ctx, "non-existent-pod")
			assert.Errorf(t, err, "pod non-existent-pod not found")
		})
	})
}

func TestPodRegistry_CreatePod(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewPodRegistry(etcdStorage)
		ctx := context.Background()

		// Test Create
		pod := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name: "test-pod",
			},
			Spec: api.PodSpec{
				Containers: []api.Container{
					{
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
			Status: api.PodPending,
		}

		err := registry.CreatePod(ctx, pod)
		require.NoError(t, err)

		// Verify pod was created
		_, err = registry.GetPod(ctx, "test-pod")
		require.NoError(t, err)
	})
}

func TestPodRegistry_UpdatePod(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewPodRegistry(etcdStorage)
		ctx := context.Background()

		pod := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name: "test-pod",
			},
			Spec: api.PodSpec{
				Containers: []api.Container{
					{
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
			Status: api.PodPending,
		}

		err := registry.CreatePod(ctx, pod)
		require.NoError(t, err)

		// Update pod status
		pod.Status = api.PodRunning
		err = registry.UpdatePod(ctx, pod)
		require.NoError(t, err)

		// Verify updated status
		retrievedPod, err := registry.GetPod(ctx, "test-pod")
		require.NoError(t, err)
		assert.Equal(t, api.PodRunning, retrievedPod.Status)
	})
}

func TestPodRegistry_DeletePod(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewPodRegistry(etcdStorage)
		ctx := context.Background()

		pod := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name: "test-pod",
			},
			Spec: api.PodSpec{
				Containers: []api.Container{
					{
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
			Status: api.PodPending,
		}

		err := registry.CreatePod(ctx, pod)
		require.NoError(t, err)

		err = registry.DeletePod(ctx, "test-pod")
		require.NoError(t, err)

		_, err = registry.GetPod(ctx, "test-pod")
		assert.Error(t, err)
	})
}

func TestPodRegistry_ListPods(t *testing.T) {
	storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
		etcdStorage := storage.NewEtcdStorage(etcdServer)
		registry := NewPodRegistry(etcdStorage)
		ctx := context.Background()

		// Test cases

		// Test ListPods
		pod1 := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name: "test-pod-1",
			},
			Spec: api.PodSpec{
				Containers: []api.Container{
					{
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
			Status: api.PodPending,
		}

		pod2 := &api.Pod{
			ObjectMeta: api.ObjectMeta{
				Name: "test-pod-2",
			},
			Spec: api.PodSpec{
				Containers: []api.Container{
					{
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
			Status: api.PodRunning,
		}

		err := registry.CreatePod(ctx, pod1)
		require.NoError(t, err)

		err = registry.CreatePod(ctx, pod2)
		require.NoError(t, err)

		pods, err := registry.ListPods(ctx)
		require.NoError(t, err)
		require.Len(t, pods, 2)

		// Verify pod names
		assert.Equal(t, "test-pod-1", pods[0].Name)
		assert.Equal(t, "test-pod-2", pods[1].Name)
	})
}

func TestPodRegistry_ListPendingPods(t *testing.T) {
	testCases := []struct {
		name                string
		podsToCreate        []*api.Pod
		expectedPendingPods int
	}{
		{
			name: "no pending pods",
			podsToCreate: []*api.Pod{
				{ObjectMeta: api.ObjectMeta{Name: "pod1"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodRunning},
				{ObjectMeta: api.ObjectMeta{Name: "pod2"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodRunning},
			},
			expectedPendingPods: 0,
		},
		{
			name: "some pending pods",
			podsToCreate: []*api.Pod{
				{ObjectMeta: api.ObjectMeta{Name: "pod3"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
				{ObjectMeta: api.ObjectMeta{Name: "pod4"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodRunning},
				{ObjectMeta: api.ObjectMeta{Name: "pod5"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
			},
			expectedPendingPods: 2,
		},
		{
			name: "all pending pods",
			podsToCreate: []*api.Pod{
				{ObjectMeta: api.ObjectMeta{Name: "pod6"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
				{ObjectMeta: api.ObjectMeta{Name: "pod7"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
			},
			expectedPendingPods: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
				etcdStorage := storage.NewEtcdStorage(etcdServer)
				registry := NewPodRegistry(etcdStorage)
				ctx := context.Background()

				// Create test pods
				for _, pod := range tc.podsToCreate {
					if err := registry.CreatePod(ctx, pod); err != nil {
						t.Fatalf("Failed to create test pod: %v", err)
					}
				}

				// Call ListPods
				pods, err := registry.ListPendingPods(ctx)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedPendingPods, len(pods))
			})
		})
	}
}

func TestPodRegistry_ListUnassignedPods(t *testing.T) {
	testCases := []struct {
		name                   string
		podsToCreate           []*api.Pod
		expectedUnassignedPods int
	}{
		{
			name: "no unassigned pods",
			podsToCreate: []*api.Pod{
				{ObjectMeta: api.ObjectMeta{Name: "pod1"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodRunning},
				{ObjectMeta: api.ObjectMeta{Name: "pod2"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodRunning},
			},
			expectedUnassignedPods: 0,
		},
		{
			name: "some unassigned pods",
			podsToCreate: []*api.Pod{
				{ObjectMeta: api.ObjectMeta{Name: "pod3"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
				{ObjectMeta: api.ObjectMeta{Name: "pod4"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodRunning},
				{ObjectMeta: api.ObjectMeta{Name: "pod5"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
			},
			expectedUnassignedPods: 2,
		},
		{
			name: "all unassigned pods",
			podsToCreate: []*api.Pod{
				{ObjectMeta: api.ObjectMeta{Name: "pod6"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
				{ObjectMeta: api.ObjectMeta{Name: "pod7"},
					Spec:   api.PodSpec{Containers: []api.Container{{Name: "test-container2", Image: "nginx"}}},
					Status: api.PodPending},
			},
			expectedUnassignedPods: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage.TestWithEmbeddedEtcd(t, func(t *testing.T, etcdServer *clientv3.Client) {
				etcdStorage := storage.NewEtcdStorage(etcdServer)
				registry := NewPodRegistry(etcdStorage)
				ctx := context.Background()

				// Create test pods
				for _, pod := range tc.podsToCreate {
					if err := registry.CreatePod(ctx, pod); err != nil {
						t.Fatalf("Failed to create test pod: %v", err)
					}
				}

				// Call ListPods
				pods, err := registry.ListUnassignedPods(ctx)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedUnassignedPods, len(pods))
			})
		})
	}
}
