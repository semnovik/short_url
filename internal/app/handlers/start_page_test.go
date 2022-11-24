package handlers

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"short_url/internal/app/repository"
	"short_url/internal/app/service"
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

			h := Handler{Server: service.NewServer(repository.NewRepository([]string{}))}
			h.startPage(w, req)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, test.want.statusCode, resp.StatusCode)
			require.Equal(t, test.want.response, string(respBody))
			require.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))

		})
	}
}
