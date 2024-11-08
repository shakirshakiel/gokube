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
