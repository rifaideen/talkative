package talkative_test

import (
	"talkative"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
