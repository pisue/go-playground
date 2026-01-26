package domain

import "errors"

// 1. 전송 계층과 상관없는 제네릭 에러 정의
var (
	ErrorBadRequest    = errors.New("bad_request")
	ErrNotFound        = errors.New("not_found")
	ErrInternalFailure = errors.New("internal_failure")
)

// 2. 커스텀 에러 구조체
type AppError struct {
	ServiceError error // 추상화된 에러 카테고리 (ErrBadRequest 등)
	Detail       error // 구체적인 기술 에러 (DB 에러 등)
	Message      string
}
