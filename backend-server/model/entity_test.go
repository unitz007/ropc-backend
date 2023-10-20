package model

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_CreateClient(t *testing.T) {

	t.Run("client secret should not be empty", func(t *testing.T) {
		client, _ := NewApplication(mock.Anything)

		if client.ClientSecret == "" {
			t.Errorf("client secret should not be empty")
		}
	})

	t.Run("client id should not be empty", func(t *testing.T) {
		_, err := NewApplication("")
		require.Error(t, err, "should return error because client id should not be empty")

	})
}
