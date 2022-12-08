package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"short_url/internal/repositories/mock_repositories"
	"testing"
)

func TestShorter_PostURL(t *testing.T) {

	type testCase struct {
		name     string
		sendURL  string
		sendUUID string
	}

	tests := []testCase{
		{name: "Simple", sendURL: "https://yandex.ru", sendUUID: "qwerty"},
		{name: "Empty", sendURL: "", sendUUID: "empty"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repositories.NewMockURLRepo(c)
			service := shorter{Repository: repo}

			genUUID = func() string {
				return test.sendUUID
			}

			repo.EXPECT().Get(test.sendUUID).Return("", errors.New("have no id"))
			repo.EXPECT().Add(test.sendUUID, test.sendURL)

			gotUUID := service.PostURL(test.sendURL)
			require.Equal(t, test.sendUUID, gotUUID)
		})
	}
}

func TestShorter_GetURLByID(t *testing.T) {

	type testCase struct {
		name     string
		sendUUID string
		wantURL  string
		err      error
	}

	tests := []testCase{
		{name: "Simple", sendUUID: "qwerty", wantURL: "https://yandex.ru", err: nil},
		{name: "doesn't exist", sendUUID: "nothing", wantURL: "", err: errors.New("url with that id is not found")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repositories.NewMockURLRepo(c)
			service := shorter{Repository: repo}

			repo.EXPECT().Get(test.sendUUID).Return(test.wantURL, test.err)

			gotURL, gotErr := service.GetURLByID(test.sendUUID)
			require.Equal(t, test.wantURL, gotURL)
			require.Equal(t, test.err, gotErr)
		})
	}
}
