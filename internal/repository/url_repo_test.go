package repository

import (
	"github.com/stretchr/testify/require"
	"short_url/configs"
	"testing"
)

func TestSimpleUrlRepo(t *testing.T) {
	configs.Config.FileStoragePath = ""
	repo := NewURLRepository()

	URL := "http://google.com"

	genUUID = func() string {
		return "qwerty"
	}
	_, err := repo.Add(URL)
	require.NoError(t, err)

	goURL, err := repo.Get("qwerty")
	require.NoError(t, err)
	require.Equal(t, URL, goURL)
}

func TestGetFromEmptyUrlRepo(t *testing.T) {
	configs.Config.FileStoragePath = ""

	repo := NewURLRepository()

	goURL, err := repo.Get("qwerty")
	require.Error(t, err)
	require.Equal(t, "", goURL)
}
