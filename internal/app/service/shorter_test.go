package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"short_url/internal/app/repository"
	mock_repository "short_url/internal/app/repository/mocks"
	"testing"
)

func TestShorter_PostURL(t *testing.T) {
	type mockBehavior func(r *mock_repository.MockURLRepo, url string)

	type testCase struct {
		name string
		send string
		want string
		repo mockBehavior
	}

	tests := []testCase{
		{
			name: "Simple positive",
			send: "http://google.com",
			want: "123",
			repo: func(r *mock_repository.MockURLRepo, url string) {
				r.EXPECT().Add(url).Return(123)
			},
		},
		{
			name: "Empty id",
			send: "",
			want: "1",
			repo: func(r *mock_repository.MockURLRepo, url string) {
				r.EXPECT().Add(url).Return(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			urlRepo := mock_repository.NewMockURLRepo(c)
			test.repo(urlRepo, test.send)

			repo := &repository.Repository{URLRepo: urlRepo}
			service := NewShorter(repo)

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
		{name: "symbols in id", send: "qwert", want: "", wantErr: true, repo: repository.NewRepository([]string{"http://google.com"})},
		{name: "id isn't set", send: "", want: "", wantErr: true, repo: repository.NewRepository([]string{"http://google.com"})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := Service{NewServer(test.repo)}
			got, err := service.GetURLByID(test.send)

			if test.wantErr == true {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.want, got)
		})
	}
}
