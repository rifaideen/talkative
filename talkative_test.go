package talkative

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("client-validation-error", func(t *testing.T) {
		client, err := New(" ")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUrl)
		assert.Nil(t, client)
	})

	t.Run("client-success", func(t *testing.T) {
		client, err := New("http://localhost:11434")

		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Assert client urls  and client are set
		assert.Equal(t, client.urls["chat"], "http://localhost:11434/api/chat")
		assert.NotNil(t, client.client)
	})
}
