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
	"short_url/pkg"
	"testing"
)

func TestShorterSrv_Shorten_HappyPass(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLStorage(ctrl)
	server := NewShorterSrv(repo)

	repository.GenUUID = func() string {
		return "328226"
	}

	repo.EXPECT().AddByUser("328226", "https://github.com/semnovik/").Return("shortUUID", nil)

	requestShoeten, err := json.Marshal(&pkg.RequestShorten{URL: "https://github.com/semnovik/"})
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestShoeten))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respPayload := new(pkg.ResponseShorten)
	err = json.NewDecoder(resp.Body).Decode(respPayload)
	require.NoError(t, err)

	require.Equal(t, configs.Config.BaseURL+"/"+"shortUUID", respPayload.Result)
}

func TestShorterSrv_Shorten_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLStorage(ctrl)
	server := NewShorterSrv(repo)

	repository.GenUUID = func() string {
		return "328226"
	}

	repo.EXPECT().AddByUser("328226", "https://github.com/semnovik/").Return("shortUUID", errors.New(`already exists`))

	requestShoeten, err := json.Marshal(&pkg.RequestShorten{URL: "https://github.com/semnovik/"})
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestShoeten))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)

	respPayload := new(pkg.ResponseShorten)
	err = json.NewDecoder(resp.Body).Decode(respPayload)
	require.NoError(t, err)

	require.Equal(t, configs.Config.BaseURL+"/"+"shortUUID", respPayload.Result)
}

func TestShorterSrv_SendURL_HappyPass(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLStorage(ctrl)
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

	repo := mock_repository.NewMockURLStorage(ctrl)
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

func TestShorterSrv_Batch_HappyPass(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLStorage(ctrl)
	server := NewShorterSrv(repo)

	repo.EXPECT().Add("https://second.com").Return("firstUUID", nil)
	repo.EXPECT().Add("https://first.com").Return("secondUUID", nil)

	requestShortenBatch, err := json.Marshal(&[]pkg.RequestShortenBatch{
		{"first", "https://second.com"},
		{"second", "https://first.com"},
	})
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(requestShortenBatch))

	server.Handler.ServeHTTP(rw, req)

	resp := rw.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var respPayload []pkg.ResponseShortenBatch
	err = json.NewDecoder(resp.Body).Decode(&respPayload)
	require.NoError(t, err)

	require.Equal(t, 2, len(respPayload))
	require.Equal(t, "first", respPayload[0].CorrelationID)
	require.Equal(t, configs.Config.BaseURL+"/"+"firstUUID", respPayload[0].ShortURL)

	require.Equal(t, "second", respPayload[1].CorrelationID)
	require.Equal(t, configs.Config.BaseURL+"/"+"secondUUID", respPayload[1].ShortURL)
}

func TestShorterSrv_GetFullURL(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	repo := mock_repository.NewMockURLStorage(ctrl)
	server := NewShorterSrv(repo)

	repo.EXPECT().Get("someUUID").Return("https://google.com", false, nil)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/someUUID", nil)

	server.Handler.ServeHTTP(rw, req)

	respPayload := rw.Result()
	response, err := io.ReadAll(respPayload.Body)
	require.NoError(t, err)
	defer respPayload.Body.Close()

	locationHeader := respPayload.Header.Get("Location")

	require.Equal(t, http.StatusTemporaryRedirect, respPayload.StatusCode)
	require.Equal(t, "https://google.com", locationHeader)
	require.Equal(t, "", string(response))
}
