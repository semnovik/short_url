package repository

import (
	"github.com/stretchr/testify/require"
	"short_url/configs"
	"testing"
)

func TestGetFromEmptyUrlRepo(t *testing.T) {
	configs.Config.FileStoragePath = ""

	repo := NewRepo(nil)

	goURL, isDeleted, err := repo.Get("qwerty")
	require.False(t, isDeleted)
	require.Error(t, err)
	require.Equal(t, "", goURL)
}
