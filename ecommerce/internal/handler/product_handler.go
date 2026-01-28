package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Create 상품 생성 핸들러
func Create(c echo.Context) error {
	req := new(CreateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	product := map[string]interface{}{
		"id":          1, // 실제로는 DB에서 생성된 ID를 사용
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
	}

	return c.JSON(http.StatusCreated, product)
}

// Get 상품 조회 핸들러
func Get(c echo.Context) error {
	id := c.Param("id")

	product := map[string]interface{}{
		"id":          id,
		"name":        "Sample Product",
		"description": "This is a sample product.",
		"price":       99.99,
	}

	return c.JSON(http.StatusOK, product)
}

// Update 상품 수정 핸들러
func Update(c echo.Context) error {
	id := c.Param("id")
	req := new(UpdateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	product := map[string]interface{}{
		"id":          id,
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
	}

	return c.JSON(http.StatusOK, product)
}

// Delete 상품 삭제 핸들러
func Delete(c echo.Context) error {
	id := c.Param("id")

	// 실제로는 DB에서 상품을 삭제하는 로직이 필요

	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted", "id": id})
}
