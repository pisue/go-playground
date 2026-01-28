package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pisue/go-playground/ecommerce/internal/handler"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Ecommerce Service!")
	})

	// 상품 관련 API 그룹
	products := e.Group("/products")
	{
		products.GET("", handler.Get)           // 상품 목록 조회
		products.POST("", handler.Create)       // 상품 생성
		products.GET("/:id", handler.Get)       // 상품 상세 조회
		products.PUT("/:id", handler.Update)    // 상품 수정
		products.DELETE("/:id", handler.Delete) // 상품 삭제
	}
	// Start server on port 8081
	e.Logger.Fatal(e.Start(":8081"))
}
