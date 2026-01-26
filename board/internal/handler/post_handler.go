package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pisue/go-playground/board/internal/domain"
	"github.com/pisue/go-playground/board/internal/service"
)

// PostHandler 구조체 정의
type PostHandler struct {
	svc service.PostService
}

// NewPostHandler 생성자 함수
func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{
		svc: svc,
	}
}

// CreateRequest 게시글 생성 요청 DTO
type CreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// UpdateRequest 게시글 수정 요청 DTO
type UpdateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Create 게시글 생성 핸들러
func (h *PostHandler) Create(c echo.Context) error {
	req := new(CreateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	post, err := h.svc.CreatePost(req.Title, req.Content, req.Author)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, post)
}

// Get 게시글 조회 핸들러
func (h *PostHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	post, err := h.svc.GetPost(uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, post)
}

// List 게시글 목록 조회 핸들러
func (h *PostHandler) List(c echo.Context) error {
	posts, err := h.svc.ListPosts()
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, posts)
}

// Update 게시글 수정 핸들러
func (h *PostHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	req := new(UpdateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	post, err := h.svc.UpdatePost(uint(id), req.Title, req.Content)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, post)
}

// Delete 게시글 삭제 핸들러
func (h *PostHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	if err := h.svc.DeletePost(uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// 에러 변환(Translation) 프로세스
func (h *PostHandler) handleError(c echo.Context, err error) error {
	var appErr *domain.AppError

	// 1. 서비스에서 반환된 에러가 우리가 정의한 AppError인지 확인
	if errors.As(err, &appErr) {
		// 2. 실제 에러 내용은 서버 로그에 남김
		if appErr.Detail != nil {
			c.Logger().Errorf("Detail Error: %v", appErr.Detail)
		}
		// 3. ServiceError(카테고리)에 따라 HTTP 상태 코드 매핑
		switch {
		case errors.Is(appErr.ServiceError, domain.ErrBadRequest):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": appErr.Message})
		case errors.Is(appErr.ServiceError, domain.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": appErr.Message})
		case errors.Is(appErr.ServiceError, domain.ErrInternalFailure):
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": appErr.Message})
		}
	}

	// 4. 정의되지 않은 일반 에러인 경우 500 에러 반환
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
