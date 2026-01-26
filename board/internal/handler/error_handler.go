package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pisue/go-playground/board/internal/domain"
)

// CustomHTTPErrorHandler 는 모든 핸들러에서 반환된 에러를 가로채서 처리
func CustomHTTPErrorHandler(err error, c echo.Context) {
	// 1. 우리가 정의한 AppError인지 확인
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		code := http.StatusInternalServerError

		// 2. 도메인 에러 카테고리에 따른 HTTP 상태 코드 매핑
		switch {
		case errors.Is(appErr.ServiceError, domain.ErrBadRequest):
			code = http.StatusBadRequest
		case errors.Is(appErr.ServiceError, domain.ErrNotFound):
			code = http.StatusNotFound
		case errors.Is(appErr.ServiceError, domain.ErrInternalFailure):
			code = http.StatusInternalServerError
		}

		// 3. 보안을 위해 상세 에러(Detail)는 서버 로그에만 남기고, 메시지만 클라이언트에게 전송
		if appErr.Detail != nil {
			c.Logger().Errorf("Detailed Error: %v", appErr.Detail)
		}

		c.JSON(code, map[string]string{"error": appErr.Message})
		return
	}

	// 4. Echo 자체 에러(예: 404 Route Not Found)나 정의되지 않은 일반 에러 처리
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
