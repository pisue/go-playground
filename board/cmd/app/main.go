package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pisue/go-playground/board/internal/handler"
	"github.com/pisue/go-playground/board/internal/repository"
	"github.com/pisue/go-playground/board/internal/service"
)

func main() {
	// 1. 의존성 주입 (Wiring)
	// Repository -> Service -> Handler 순으로 생성
	postRepo := repository.NewMemoryPostRepository()
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	// 2. Echo 인스턴스 설정
	e := echo.New()

	// 미들웨어 설정
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 에러 핸들러 등록
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// 3. 라우팅 설정
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Board API Server is running!")
	})

	// 게시글 관련 API 그룹
	posts := e.Group("/posts")
	{
		posts.POST("", postHandler.Create)       // 게시글 생성
		posts.GET("", postHandler.List)          // 게시글 목록 조회
		posts.GET("/:id", postHandler.Get)       // 게시글 상세 조회
		posts.PUT("/:id", postHandler.Update)    // 게시글 수정
		posts.DELETE("/:id", postHandler.Delete) // 게시글 삭제
	}

	// 4. 서버 시작
	e.Logger.Fatal(e.Start(":8080"))
}
