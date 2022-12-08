package storage

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleUrlRepo(t *testing.T) {

	repo := NewURLRepository()

	URL := "http://google.com"
	repo.Add("123", URL)

	goURL, err := repo.Get("123")
	require.NoError(t, err)
	require.Equal(t, URL, goURL)
}

func TestGetFromEmptyUrlRepo(t *testing.T) {

	repo := NewURLRepository()

	goURL, err := repo.Get("qwerty")
	require.Error(t, err)
	require.Equal(t, "", goURL)
}
