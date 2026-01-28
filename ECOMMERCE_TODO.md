# 🛒 E-commerce Service TODO List

이 문서는 E-commerce 마이크로서비스(`ecommerce`) 개발을 위한 작업 목록을 관리합니다.

## ✅ Phase 1: 기본 설정 및 에러 핸들링 (진행 중)
- [x] 프로젝트 초기화 및 Echo 설정
- [x] 도메인 에러(`AppError`) 정의 (`internal/domain`)
- [x] 글로벌 에러 핸들러(`CustomHTTPErrorHandler`) 구현
- [x] `/products` 라우팅 뼈대 작성

## 🚀 Phase 2: 상품(Product) CRUD 구현
- [ ] **Domain**: `Product` 구조체 정의 (`internal/domain/product.go`)
- [ ] **Repository (Interface)**: `ProductRepository` 인터페이스 정의
- [ ] **Repository (Impl)**: `ProductRepository` 메모리/DB 구현체 작성
- [ ] **Service**: `ProductService` 구현 (비즈니스 로직)
- [ ] **Handler**: `ProductHandler` 구현 및 `main.go` 연결
- [ ] **Test**: 핸들러 및 서비스 단위 테스트 작성

## 📦 Phase 3: 데이터베이스 연동 (추후 예정)
- [ ] 데이터베이스(MySQL/PostgreSQL) 드라이버 설정
- [ ] 마이그레이션 스크립트 작성
