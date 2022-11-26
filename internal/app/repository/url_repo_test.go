package repository

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepoURL_Add(t *testing.T) {

	type testCase struct {
		name    string
		sendUrl string
		wantId  int
		repo    URLRepo
	}

	tests := []testCase{
		{name: "add to empty", sendUrl: "http://google.com", wantId: 1, repo: NewRepository([]string{})},
		{name: "add when some exists", sendUrl: "http://google.com", wantId: 2, repo: NewRepository([]string{"some"})},
		{name: "add empty string", sendUrl: "", wantId: 1, repo: NewRepository([]string{})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotId := test.repo.Add(test.sendUrl)

			require.Equal(t, test.wantId, gotId)
		})
	}

}

func TestRepoURL_Get(t *testing.T) {

	type testCase struct {
		name    string
		sendId  int
		wantUrl string
		wantErr bool
		repo    URLRepo
	}

	tests := []testCase{
		{name: "get existing url", sendId: 1, wantUrl: "http://google.com", wantErr: false, repo: NewRepository([]string{"http://google.com"})},
		{name: "get url, which doesn't exist ", sendId: 2, wantUrl: "", wantErr: true, repo: NewRepository([]string{"http://google.com"})},
		{name: "get second existing url", sendId: 2, wantUrl: "http://google.com", wantErr: false, repo: NewRepository([]string{"ya.ru", "http://google.com"})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUrl, err := test.repo.Get(test.sendId)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.wantUrl, gotUrl)
		})
	}
}
