package service

import (
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
