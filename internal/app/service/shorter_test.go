package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"short_url/internal/app/repository"
	"testing"
)

func TestShorter_PostURL(t *testing.T) {

	type testCase struct {
		name string
		send string
		want string
		repo *repository.Repository
	}

	tests := []testCase{
		{name: "Simple positive", send: "http://google.com", want: "1", repo: repository.NewRepository([]string{})},
		{name: "Empty id", send: "", want: "1", repo: repository.NewRepository([]string{})},
		{name: "Repo is not empty", send: "http://google.com", want: "2", repo: repository.NewRepository([]string{"http://some.com"})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := Service{NewServer(test.repo)}
			want := service.PostURL(test.send)
			require.Equal(t, test.want, want)
		})
	}

}

func TestShorter_GetURLByID(t *testing.T) {

	type testCase struct {
		name    string
		send    string
		want    string
		wantErr bool
		repo    *repository.Repository
	}

	tests := []testCase{
		{name: "Simple positive", send: "1", want: "http://google.com", wantErr: false, repo: repository.NewRepository([]string{"http://google.com"})},
		{name: "More than one url in repo", send: "2", want: "http://yandex.com", wantErr: false, repo: repository.NewRepository([]string{"http://google.com", "http://yandex.com"})},
		{name: "Not found", send: "123", want: "", wantErr: true, repo: repository.NewRepository([]string{"http://google.com"})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := Service{NewServer(test.repo)}
			got, err := service.GetURLByID(test.send)

			if test.wantErr == true {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.want, got)
		})
	}
}
