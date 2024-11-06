package api

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestContainerValidation(t *testing.T) {
	validate := validator.New()

	t.Run("should validate container with required fields", func(t *testing.T) {
		container := Container{
			Name:  "nginx-container",
			Image: "nginx:latest",
		}

		err := validate.Struct(container)
		assert.NoError(t, err)
	})

	t.Run("should fail validation if name is missing", func(t *testing.T) {
		container := Container{
			Image: "nginx:latest",
		}

		err := validate.Struct(container)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Container.Name' Error:Field validation for 'Name' failed on the 'required' tag")
	})

	t.Run("should fail validation if image is missing", func(t *testing.T) {
		container := Container{
			Name: "nginx-container",
		}

		err := validate.Struct(container)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Container.Image' Error:Field validation for 'Image' failed on the 'required' tag")
	})
}

func TestPodSpecValidation(t *testing.T) {
	validate := validator.New()

	t.Run("should validate PodSpec with required fields", func(t *testing.T) {
		podSpec := PodSpec{
			Containers: []Container{
				{
					Name:  "nginx-container",
					Image: "nginx:latest",
				},
			},
			Replicas: 3,
		}

		err := validate.Struct(podSpec)
		assert.NoError(t, err)
	})

	t.Run("should fail validation if containers are missing", func(t *testing.T) {
		podSpec := PodSpec{
			Replicas: 3,
		}

		err := validate.Struct(podSpec)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'PodSpec.Containers' Error:Field validation for 'Containers' failed on the 'required' tag")
	})

	t.Run("should fail validation if container image is missing", func(t *testing.T) {
		podSpec := PodSpec{
			Containers: []Container{
				{
					Name: "nginx-container",
				},
			},
			Replicas: 3,
		}

		err := validate.Struct(podSpec)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'PodSpec.Containers[0].Image' Error:Field validation for 'Image' failed on the 'required' tag")
	})

	t.Run("should fail validation if replicas is negative", func(t *testing.T) {
		podSpec := PodSpec{
			Containers: []Container{
				{
					Name:  "nginx-container",
					Image: "nginx:latest",
				},
			},
			Replicas: -1,
		}

		err := validate.Struct(podSpec)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'PodSpec.Replicas' Error:Field validation for 'Replicas' failed on the 'gte' tag")
	})
}

func TestPodValidation(t *testing.T) {
	validate := validator.New()

	t.Run("should validate Pod with required fields", func(t *testing.T) {
		pod := Pod{
			ObjectMeta: ObjectMeta{
				Name: "test-pod",
			},
			Spec: PodSpec{
				Containers: []Container{
					{
						Name:  "nginx-container",
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
			Status: PodPending,
		}

		err := validate.Struct(pod)
		assert.NoError(t, err)
	})

	t.Run("should fail validation if spec is missing", func(t *testing.T) {
		pod := Pod{
			ObjectMeta: ObjectMeta{
				Name: "test-pod",
			},
			Status: PodPending,
		}

		err := validate.Struct(pod)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Pod.Spec.Containers' Error:Field validation for 'Containers' failed on the 'required' tag")
	})

	t.Run("should fail validation if status is missing", func(t *testing.T) {
		pod := Pod{
			ObjectMeta: ObjectMeta{
				Name: "test-pod",
			},
			Spec: PodSpec{
				Containers: []Container{
					{
						Name:  "nginx-container",
						Image: "nginx:latest",
					},
				},
				Replicas: 3,
			},
		}

		err := validate.Struct(pod)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Pod.Status' Error:Field validation for 'Status' failed on the 'required' tag")
	})
}
