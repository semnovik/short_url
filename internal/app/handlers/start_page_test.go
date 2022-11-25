package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"short_url/internal/app/repository"
	"short_url/internal/app/service"
	mock_service "short_url/internal/app/service/mocks"
	"strings"
	"testing"
)

func TestStartPagePostURL(t *testing.T) {

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			name:   "Post new url",
			method: http.MethodPost,
			path:   "/",
			body:   "https://google.com",
			want: want{
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/1",
				contentType: "",
			},
		},
		{
			name:   "Post empty url",
			method: http.MethodPost,
			path:   "/",
			body:   "",
			want: want{
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/1",
				contentType: "",
			},
		},
		{
			name:   "Wrong http method",
			method: http.MethodPut,
			path:   "/",
			body:   "https://google.com",
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "Method not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(test.body)
			req := httptest.NewRequest(test.method, test.path, body)
			w := httptest.NewRecorder()

			h := Handler{Service: service.NewServer(repository.NewRepository([]string{}))}
			h.startPage(w, req)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)

			require.Equal(t, test.want.statusCode, resp.StatusCode)
			require.Equal(t, test.want.response, string(respBody))
			require.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))

		})
	}
}

func TestHandler_StartPage(t *testing.T) {
	type mockBehavior func(s *mock_service.MockShorterService, id string)

	type want struct {
		StatusCode int
		Header     string
		ExpErr     bool
	}

	tests := []struct {
		name         string
		method       string
		path         string
		request      string
		mockBehavior mockBehavior
		want         want
	}{
		{
			name:    "happy pass",
			method:  http.MethodGet,
			path:    "/",
			request: "1",
			mockBehavior: func(s *mock_service.MockShorterService, id string) {
				s.EXPECT().GetURLByID(id).Return("http://google.com", nil)
			},
			want: want{
				StatusCode: http.StatusTemporaryRedirect,
				Header:     "http://google.com",
				ExpErr:     false,
			},
		},
		{
			name:    "error occurring",
			method:  http.MethodGet,
			path:    "/",
			request: "qwerty",
			mockBehavior: func(s *mock_service.MockShorterService, id string) {
				s.EXPECT().GetURLByID(id).Return("", errors.New("something went wrong"))
			},
			want: want{
				StatusCode: http.StatusBadRequest,
				Header:     "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			shorter := mock_service.NewMockShorterService(c)
			test.mockBehavior(shorter, test.request)

			shorterService := &service.Service{ShorterService: shorter}
			handler := NewHandler(shorterService)

			//r := http.NewServeMux()
			//r.HandleFunc(test.path, handler.startPage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, test.path+test.request, nil)

			handler.startPage(w, req)
			res := w.Result()

			assert.Equal(t, test.want.StatusCode, res.StatusCode)
			assert.Equal(t, test.want.Header, res.Header.Get("Location"))

			//if test.want.ExpErr == true {
			//	assert.NotNil(t, )
			//}
		})
	}
}
