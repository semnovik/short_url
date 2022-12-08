package repositories

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleUrlRepo(t *testing.T) {

	repo := NewURLRepo()

	URL := "http://google.com"
	UUID := repo.Add(URL)
	require.NotEqual(t, "", UUID)

	goURL, err := repo.Get(UUID)
	require.NoError(t, err)
	require.Equal(t, URL, goURL)
}

func TestGetFromEmptyUrlRepo(t *testing.T) {

	repo := NewURLRepo()

	goURL, err := repo.Get("qwerty")
	require.Error(t, err)
	require.Equal(t, "", goURL)
}
