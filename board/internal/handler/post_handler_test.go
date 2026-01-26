package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pisue/go-playground/board/internal/domain"
	"github.com/pisue/go-playground/board/internal/handler"
	"github.com/stretchr/testify/assert"
)

// --- Mock Service ---
// 테스트를 위해 Service 인터페이스를 간단하게 Mocking 합니다.
type mockPostService struct{}

func (m *mockPostService) CreatePost(title, content, author string) (*domain.Post, error) {
	return nil, nil
}

// GetPost: id가 999면 ErrNotFound 반환, 그 외엔 성공 가정
func (m *mockPostService) GetPost(id uint) (*domain.Post, error) {
	if id == 999 {
		return nil, &domain.AppError{
			ServiceError: domain.ErrNotFound,
			Message:      "게시글을 찾을 수 없습니다.",
		}
	}
	if id == 500 {
		return nil, &domain.AppError{
			ServiceError: domain.ErrInternalFailure,
			Message:      "내부 서버 오류입니다.",
		}
	}
	return &domain.Post{ID: id, Title: "Test"}, nil
}

func (m *mockPostService) ListPosts() ([]*domain.Post, error) {
	return nil, nil
}
func (m *mockPostService) UpdatePost(id uint, title, content string) (*domain.Post, error) {
	return nil, nil
}
func (m *mockPostService) DeletePost(id uint) error {
	return nil
}

// --- Test Code ---

func TestPostHandler_ErrorHandling(t *testing.T) {
	// 1. Setup Echo & Dependencies
	e := echo.New()
	// 중요: 우리가 만든 글로벌 에러 핸들러를 등록해야 테스트 가능
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	svc := &mockPostService{}
	h := handler.NewPostHandler(svc)

	// 2. Define Test Cases
	tests := []struct {
		name           string
		targetURL      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "404 Not Found Error",
			targetURL:      "/posts/999", // Mock에서 에러 발생 유도
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"게시글을 찾을 수 없습니다."`,
		},
		{
			name:           "500 Internal Server Error",
			targetURL:      "/posts/500", // Mock에서 에러 발생 유도
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"내부 서버 오류입니다."`,
		},
		{
			name:           "200 OK Success",
			targetURL:      "/posts/1",
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":1`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 3. Create Request
			req := httptest.NewRequest(http.MethodGet, tt.targetURL, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Path Param 설정 (Echo 라우팅 시뮬레이션)
			// /posts/:id 형태이므로 ParamName과 ParamValues를 수동으로 설정해줘야 함
			c.SetPath("/posts/:id")
			// URL에서 ID 추출 (/posts/999 -> 999)
			parts := strings.Split(tt.targetURL, "/")
			if len(parts) > 2 {
				c.SetParamNames("id")
				c.SetParamValues(parts[2])
			}

			// 4. Execute Handler
			// 에러가 발생하면 핸들러는 에러를 리턴하고, Echo가 HTTPErrorHandler를 호출하는 구조임.
			// 단위 테스트에서는 ServeHTTP를 통하지 않고 Handler를 직접 호출하므로,
			// 리턴된 에러를 직접 HTTPErrorHandler에 넘겨줘야 통합 동작을 검증할 수 있음.
			err := h.Get(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			// 5. Assertions
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
