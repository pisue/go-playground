# Go Playground

이 저장소는 **Java Spring 개발자의 Go 언어 친해지기 프로젝트**입니다. 
Go 언어의 특징을 익히고, Spring 프레임워크에서의 개발 경험을 Go 환경에 어떻게 녹여낼 수 있는지 실습하는 것을 목표로 합니다.

Go Workspaces(`go.work`)를 활용하여 여러 Go 프로젝트(`board`, `ecommerce`)를 하나의 워크스페이스에서 효율적으로 관리합니다.

## 왜 Go Workspaces를 사용하나요?

1.  **멀티 모듈 관리**: 하나의 저장소 내에서 독립된 여러 모듈(`board`, `ecommerce` 등)을 동시에 개발하기 편리합니다.
2.  **로컬 의존성 해결**: 모듈 간의 의존성이 있을 때 `go.mod` 파일에 일일이 `replace` 구문을 추가하지 않고도 로컬에서 유연하게 참조할 수 있습니다.
3.  **Spring 멀티 프로젝트와 유사한 경험**: Java/Spring의 멀티 모듈 프로젝트(Maven/Gradle) 구조와 유사한 관리 환경을 제공하여 익숙한 방식으로 프로젝트를 구성할 수 있습니다.

## 프로젝트 구성

- **board**: Echo 프레임워크 기반의 계층형 아키텍처(Layered Architecture) 게시판 프로젝트
- **ecommerce**: 향후 추가 예정인 이커머스 프로젝트

## 디렉토리 구조 (Board)

[Standard Go Project Layout](./ARCHITECTURE.md)을 따릅니다.
[(이 Layout은 Go 스러운가?)](./GO_IDIOMATIC_LAYOUT.md)

```text
board/
├── cmd/app/main.go        # 애플리케이션 진입점
├── internal/
│   ├── handler/           # HTTP 요청 처리 (Controller)
│   ├── service/           # 비즈니스 로직
│   ├── repository/        # 데이터베이스 접근
│   └── model/             # 구조체 정의 (Entity/DTO)
├── pkg/utils/             # 공용 유틸리티
└── configs/               # 설정 파일
```

## 초기화 및 실행 방법

### 1. 워크스페이스 및 모듈 초기화

시스템에 Go가 설치되어 있어야 합니다. 루트 디렉토리(`go-playground`)에서 다음 명령을 실행하세요.

```bash
# 워크스페이스 초기화
go work init ./board ./ecommerce

# board 모듈 설정
cd board
go mod init github.com/pisue/go-playground/board
go get github.com/labstack/echo/v4

# ecommerce 모듈 설정
cd ../ecommerce
go mod init github.com/pisue/go-playground/ecommerce
```

### 2. 서버 실행

```bash
cd board
go run cmd/app/main.go
```

서버가 실행되면 `http://localhost:8080`에서 "Hello, World!" 메시지를 확인할 수 있습니다.
