package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"short_url/internal/app/repository"
	mock_repository "short_url/internal/app/repository/mocks"
	"strconv"
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
	type mockBehavior func(r *mock_repository.MockURLRepo, urlId int)

	type testCase struct {
		name    string
		send    int
		want    string
		wantErr bool
		repo    mockBehavior
	}

	tests := []testCase{
		{
			name:    "Simple positive",
			send:    1,
			want:    "http://google.com",
			wantErr: false,
			repo: func(r *mock_repository.MockURLRepo, urlId int) {
				r.EXPECT().Get(urlId).Return("http://google.com", nil)
			},
		},
		{
			name:    "More than one url in repo",
			send:    2,
			want:    "http://yandex.com",
			wantErr: false,
			repo: func(r *mock_repository.MockURLRepo, urlId int) {
				r.EXPECT().Get(urlId).Return("http://yandex.com", nil)
			},
		},
		{
			name:    "Not found",
			send:    123,
			want:    "",
			wantErr: true,
			repo: func(r *mock_repository.MockURLRepo, urlId int) {
				r.EXPECT().Get(urlId).Return("", errors.New("url with that id is not found"))
			},
		},
		{
			name:    "id is zero",
			send:    0,
			want:    "",
			wantErr: true,
			repo: func(r *mock_repository.MockURLRepo, urlId int) {
				r.EXPECT().Get(urlId).Return("", errors.New("url with that id is not found"))
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

			service := &Service{NewShorter(repo)}
			got, err := service.GetURLByID(strconv.Itoa(test.send))

			if test.wantErr == true {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.want, got)
		})
	}
}
