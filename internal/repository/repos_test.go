package repository

//
//import (
//	"github.com/stretchr/testify/require"
//	"short_url/configs"
//	"testing"
//)
//
//func TestSimpleUrlRepo(t *testing.T) {
//	configs.Config.FileStoragePath = ""
//	repo := NewRepo()
//
//	URL := "http://google.com"
//
//	GenUUID = func() string {
//		return "qwerty"
//	}
//	_, err := repo.Add(URL)
//	require.NoError(t, err)
//
//	goURL, err := repo.Get("qwerty")
//	require.NoError(t, err)
//	require.Equal(t, URL, goURL)
//}
//
//func TestGetFromEmptyUrlRepo(t *testing.T) {
//	configs.Config.FileStoragePath = ""
//
//	repo := NewRepo()
//
//	goURL, err := repo.Get("qwerty")
//	require.Error(t, err)
//	require.Equal(t, "", goURL)
//}

//func TestPostgresRepo_Add(t *testing.T) {
//	configs.Config.DatabaseDSN = "Not empty"
//
//	repo := mock_repository.MockURLRepo{}
//}
