package talkative_test

import (
	"talkative"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient tests the behavior of the New function in two scenarios:
// 1. When the URL parameter is empty or contains only whitespace characters, it should return an error of type talkative.ErrUrl and a nil client.
// 2. When a valid URL is provided, it should return a non-nil client and no error.
//
// Parameters:
// - t: A testing.T object used for running the test and reporting any failures.
func TestClient(t *testing.T) {
	t.Run("client-validation-error", func(t *testing.T) {
		client, err := talkative.New(" ")

		assert.Error(t, err)
		assert.ErrorIs(t, err, talkative.ErrUrl)
		assert.Nil(t, client)
	})

	t.Run("client-success", func(t *testing.T) {
		client, err := talkative.New("http://localhost:11434")

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
