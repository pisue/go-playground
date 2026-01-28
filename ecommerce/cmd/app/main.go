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

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Ecommerce Service!")
	})

	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// 상품 관련 API 그룹
	products := e.Group("/products")
	{
		products.GET("", func(c echo.Context) error {
			return c.String(http.StatusOK, "List of products")
		})
		products.POST("", func(c echo.Context) error {
			return c.String(http.StatusCreated, "Product created")
		})
		products.GET("/:id", func(c echo.Context) error {
			return c.String(http.StatusOK, "Product details for ID: "+c.Param("id"))
		})
		products.PUT("/:id", func(c echo.Context) error {
			return c.String(http.StatusOK, "Product updated for ID: "+c.Param("id"))
		})
		products.DELETE("/:id", func(c echo.Context) error {
			return c.String(http.StatusOK, "Product deleted for ID: "+c.Param("id"))
		})
	}
	// Start server on port 8081
	e.Logger.Fatal(e.Start(":8081"))
}
