package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pisue/go-playground/ecommerce/internal/domain"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		code := http.StatusInternalServerError

		switch {
		case errors.Is(appErr.ServiceError, domain.ErrBadRequest):
			code = http.StatusBadRequest
		case errors.Is(appErr.ServiceError, domain.ErrNotFound):
			code = http.StatusNotFound
		case errors.Is(appErr.ServiceError, domain.ErrInternalFailure):
			code = http.StatusInternalServerError
		}

		if appErr.Detail != nil {
			c.Logger().Errorf("Detailed Error: %v", appErr.Detail)
		}

		c.JSON(code, map[string]string{"error": appErr.Message})
		return
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
