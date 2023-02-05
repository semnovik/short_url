package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"short_url/configs"
	"short_url/internal/repository"
	mock_repository "short_url/internal/repository/mock"
	"testing"
)

func TestShorterSrv_Shorten_HappyPass(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLRepo(ctrl)
	server := NewShorterSrv(repo)

	repository.GenUUID = func() string {
		return "328226"
	}

	repo.EXPECT().AddByUser("328226", "https://github.com/semnovik/").Return("shortUUID", nil)

	requestShoeten, err := json.Marshal(&RequestShorten{URL: "https://github.com/semnovik/"})
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestShoeten))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respPayload := new(ResponseShorten)
	err = json.NewDecoder(resp.Body).Decode(respPayload)
	require.NoError(t, err)

	require.Equal(t, configs.Config.BaseURL+"/"+"shortUUID", respPayload.Result)
}

func TestShorterSrv_Shorten_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLRepo(ctrl)
	server := NewShorterSrv(repo)

	repository.GenUUID = func() string {
		return "328226"
	}

	repo.EXPECT().AddByUser("328226", "https://github.com/semnovik/").Return("shortUUID", errors.New(`already exists`))

	requestShoeten, err := json.Marshal(&RequestShorten{URL: "https://github.com/semnovik/"})
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestShoeten))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)

	respPayload := new(ResponseShorten)
	err = json.NewDecoder(resp.Body).Decode(respPayload)
	require.NoError(t, err)

	require.Equal(t, configs.Config.BaseURL+"/"+"shortUUID", respPayload.Result)
}

func TestShorterSrv_SendURL_HappyPass(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLRepo(ctrl)
	server := NewShorterSrv(repo)

	repository.GenUUID = func() string {
		return "328226"
	}

	repo.EXPECT().AddByUser("328226", "https://github.com/semnovik/").Return("shortUUID", nil)

	requestShoeten := "https://github.com/semnovik/"

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(requestShoeten)))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	resBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.Equal(t, configs.Config.BaseURL+"/"+"shortUUID", string(resBody))
}

func TestShorterSrv_SendURL_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLRepo(ctrl)
	server := NewShorterSrv(repo)

	repository.GenUUID = func() string {
		return "328226"
	}

	repo.EXPECT().AddByUser("328226", "https://github.com/semnovik/").Return("shortUUID", errors.New(`already exist`))

	requestShoeten := "https://github.com/semnovik/"

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(requestShoeten)))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	resBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
	require.Equal(t, configs.Config.BaseURL+"/"+"shortUUID", string(resBody))
}
