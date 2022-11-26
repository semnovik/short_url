package repository

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepoURL_Add(t *testing.T) {

	type testCase struct {
		name    string
		sendURL string
		wantID  int
		repo    URLRepo
	}

	tests := []testCase{
		{name: "add to empty", sendURL: "http://google.com", wantID: 1, repo: NewRepository([]string{})},
		{name: "add when some exists", sendURL: "http://google.com", wantID: 2, repo: NewRepository([]string{"some"})},
		{name: "add empty string", sendURL: "", wantID: 1, repo: NewRepository([]string{})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotID := test.repo.Add(test.sendURL)

			require.Equal(t, test.wantID, gotID)
		})
	}

}

func TestRepoURL_Get(t *testing.T) {

	type testCase struct {
		name    string
		sendID  int
		wantURL string
		wantErr bool
		repo    URLRepo
	}

	tests := []testCase{
		{name: "get existing url", sendID: 1, wantURL: "http://google.com", wantErr: false, repo: NewRepository([]string{"http://google.com"})},
		{name: "get url, which doesn't exist ", sendID: 2, wantURL: "", wantErr: true, repo: NewRepository([]string{"http://google.com"})},
		{name: "get second existing url", sendID: 2, wantURL: "http://google.com", wantErr: false, repo: NewRepository([]string{"ya.ru", "http://google.com"})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUrl, err := test.repo.Get(test.sendID)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.wantURL, gotUrl)
		})
	}
}
