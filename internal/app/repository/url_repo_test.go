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
		{name: "add to empty", sendUrl: "http://google.com", wantId: 1, repo: NewURLRepo([]string{})},
		{name: "add when some exists", sendUrl: "http://google.com", wantId: 2, repo: NewURLRepo([]string{"some"})},
		{name: "add empty string", sendUrl: "", wantId: 1, repo: NewURLRepo([]string{})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotId := test.repo.Add(test.sendUrl)

			require.Equal(t, test.wantId, gotId)
		})
	}

}
