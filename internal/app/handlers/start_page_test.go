package handlers

//
//import (
//	"errors"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"short_url/internal/app/services"
//	mock_service "short_url/internal/app/services/mocks"
//	"strings"
//	"testing"
//)
//
//func TestHandler_GetFullURL(t *testing.T) {
//	type mockBehavior func(s *mock_service.MockShorterService, id string)
//
//	type want struct {
//		StatusCode int
//		Header     string
//	}
//
//	tests := []struct {
//		name         string
//		method       string
//		path         string
//		request      string
//		mockBehavior mockBehavior
//		want         want
//	}{
//		{
//			name:    "happy pass",
//			method:  http.MethodGet,
//			path:    "/1",
//			request: "",
//			mockBehavior: func(s *mock_service.MockShorterService, id string) {
//				s.EXPECT().GetURLByID(id).Return("http://google.com", nil)
//			},
//			want: want{
//				StatusCode: http.StatusTemporaryRedirect,
//				Header:     "http://google.com",
//			},
//		},
//		{
//			name:    "error occurring",
//			method:  http.MethodGet,
//			path:    "/qwerty\"",
//			request: "",
//			mockBehavior: func(s *mock_service.MockShorterService, id string) {
//				s.EXPECT().GetURLByID(id).Return("", errors.New("something went wrong"))
//			},
//			want: want{
//				StatusCode: http.StatusBadRequest,
//				Header:     "",
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//
//			// Инициализация контроллера gomock
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			shorter := mock_service.NewMockShorterService(c)
//			test.mockBehavior(shorter, test.request)
//
//			// Инициализация слоя services с моком ShorterService
//			shorterService := &services.Service{ShorterService: shorter}
//			handler := NewHandler(shorterService)
//
//			// Инициализация тестового клиента w и запроса req
//			w := httptest.NewRecorder()
//			req := httptest.NewRequest(test.method, test.path+test.request, nil)
//
//			// Выполнение запроса и получение результатов
//			handler.GetFullURL(w, req)
//			res := w.Result()
//			defer res.Body.Close()
//
//			// Сравнение фактических результатов с ожидаемыми
//			assert.Equal(t, test.want.StatusCode, res.StatusCode)
//			assert.Equal(t, test.want.Header, res.Header.Get("Location"))
//		})
//	}
//}
//
//func TestHandler_SendURL(t *testing.T) {
//	type mockBehavior func(s *mock_service.MockShorterService, url string)
//
//	type want struct {
//		StatusCode int
//		ExpErr     bool
//		Response   string
//	}
//
//	tests := []struct {
//		name         string
//		method       string
//		path         string
//		requestBody  string
//		mockBehavior mockBehavior
//		want         want
//	}{
//		{
//			name:        "happy pass",
//			method:      http.MethodPost,
//			path:        "/",
//			requestBody: "http://google.com",
//			mockBehavior: func(s *mock_service.MockShorterService, url string) {
//				s.EXPECT().PostURL(url).Return("1")
//			},
//			want: want{
//				StatusCode: http.StatusCreated,
//				ExpErr:     false,
//				Response:   "http://localhost:8080/1",
//			},
//		},
//		{
//			name:        "not first in repo",
//			method:      http.MethodPost,
//			path:        "/",
//			requestBody: "http://google.com",
//			mockBehavior: func(s *mock_service.MockShorterService, url string) {
//				s.EXPECT().PostURL(url).Return("328225")
//			},
//			want: want{
//				StatusCode: http.StatusCreated,
//				ExpErr:     false,
//				Response:   "http://localhost:8080/328225",
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//
//			// Инициализация контроллера gomock
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			shorter := mock_service.NewMockShorterService(c)
//			test.mockBehavior(shorter, test.requestBody)
//
//			// Инициализация слоя services с моком ShorterService
//			shorterService := &services.Service{ShorterService: shorter}
//			handler := NewHandler(shorterService)
//
//			// Инициализация тестового клиента w и запроса req
//			w := httptest.NewRecorder()
//
//			req := httptest.NewRequest(test.method, test.path, strings.NewReader(test.requestBody))
//
//			// Выполнение запроса и получение результатов
//			handler.SendURL(w, req)
//			res := w.Result()
//
//			resBody, err := io.ReadAll(res.Body)
//			defer res.Body.Close()
//
//			// Сравнение фактических результатов с ожидаемыми
//			assert.NoError(t, err)
//			assert.Equal(t, test.want.StatusCode, res.StatusCode)
//			assert.Equal(t, test.want.Response, string(resBody))
//		})
//	}
//}
