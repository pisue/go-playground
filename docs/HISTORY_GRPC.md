# gRPC Server 학습 및 개발 이력

이 문서는 `grpc-server` 프로젝트를 진행하며 학습한 내용과 개발 변경 사항을 기록합니다.
기존에 진행했던 **Echo 프레임워크 기반의 HTTP 웹 서버** 프로젝트와 달리, 본 프로젝트는 **gRPC 프로토콜의 핵심 원리와 동작 방식**을 깊이 있게 이해하는 것을 최우선 목표로 합니다.

## 1. 프로젝트 개요 및 목적

*   **목표**: gRPC, Protocol Buffers(Protobuf)의 이해 및 Go 언어에서의 gRPC 서버 구현 학습.
*   **차별점**: 프레임워크의 편리한 기능(Magic)에 의존하기보다, **Pure Go**에 가까운 방식으로 직접 의존성을 주입하고 서버를 구성하여 내부 동작 흐름을 명확히 파악합니다.

## 2. 기존 프로젝트(Board/Ecommerce)와의 아키텍처 및 구현 비교

본 섹션에서는 Echo 프레임워크 기반의 `board`, `ecommerce` 모듈과 현재 `grpc-server` 모듈의 아키텍처 및 의존성 주입(DI) 방식의 차이를 분석합니다.

### 2.1 통신 프로토콜 및 진입점 (Entry Point)

| 구분 | Echo 프로젝트 (`board`) | gRPC 프로젝트 (`grpc-server`) |
| :--- | :--- | :--- |
| **기반 프로토콜** | HTTP/1.1 (JSON Text) | HTTP/2 (Protobuf Binary) |
| **라우팅** | `e.GET("/posts", handler.List)` (URL Path) | `.proto` 서비스 메서드 호출 (RPC) |
| **데이터 바인딩** | `c.Bind(&dto)` (Reflection) | `protoc` 생성 코드 (Static Typing) |

### 2.2 의존성 주입(DI) 및 포인터 사용 전략 비교

두 프로젝트 모두 `Main` 함수에서 의존성을 조립(Wiring)하지만, **추상화 수준과 포인터 노출 여부**에서 큰 차이를 보입니다.

#### A. Board 모듈: 인터페이스 기반 (Loose Coupling)
`board` 모듈은 **철저한 인터페이스 기반 설계**를 따릅니다.

*   **코드 형태**:
    ```go
    // board/cmd/app/main.go
    // 반환 타입이 Interface (구체적인 구조체 *memoryPostRepository가 아님)
    postRepo := repository.NewMemoryPostRepository() 
    postSvc := service.NewPostService(postRepo)
    ```
*   **포인터 사용 (Implicit)**:
    *   코드 상에서 `*` (Asterisk)가 보이지 않지만, 내부적으로 인터페이스 변수는 구현체의 **포인터**를 담고 있습니다.
    *   `NewPostService`의 리턴 타입은 `PostService` (Interface)입니다.
*   **장점**:
    *   **느슨한 결합**: `main.go`는 `memoryPostRepository`인지 `mysqlPostRepository`인지 알 필요가 없습니다.
    *   **테스트 용이성**: `postSvc`에 가짜 객체(Mock)를 주입하기 매우 쉽습니다.

#### B. gRPC Server 모듈: 구조체 포인터 기반 (Tight Coupling & Explicit)
현재 `grpc-server` 모듈은 **구체적인 구조체 포인터**를 직접 주고받습니다.

*   **코드 형태**:
    ```go
    // grpc-server/cmd/app.go
    // 반환 타입이 *Repository (구체적인 구조체 포인터)
    if a.repository, err = repository.NewRepository(cfg); err != nil { ... }
    // 인자로 *Repository 포인터를 직접 요구
    if a.service, err = service.NewService(cfg, a.repository); err != nil { ... }
    ```
*   **포인터 사용 (Explicit)**:
    *   `NewService` 함수 시그니처를 보면 `*repository.Repository`를 인자로 받고, `*Service`를 반환합니다.
    *   포인터 사용이 명시적으로 드러나 있습니다.
*   **의도 및 특징**:
    *   초기 개발 단계에서 **구조와 데이터 흐름을 명확히** 하기 위해 직관적인 방식을 택했습니다.
    *   하지만 이는 **강한 결합**을 초래하므로, 추후 테스트 코드 작성이나 저장소 교체 시 유연성이 떨어질 수 있습니다. (리팩토링 대상)

### 2.3 요약

*   **Board**: "나는 `PostRepository`라는 **행위(Interface)**가 필요해. 누가 오든 상관없어." (유연함)
*   **gRPC Server**: "나는 `*repository.Repository`라는 **실체(Struct Pointer)**가 필요해. 정확히 얘여야만 해." (명확하지만 경직됨)

이러한 차이를 이해하고, 향후 `grpc-server`도 인터페이스 기반으로 리팩토링하여 유연성을 확보하는 과정을 학습할 예정입니다.

## 3. 개발 진행 로그

### 3.1 초기 스캐폴딩 (Layered Architecture)
*   **Flag 기반 설정 관리**:
    *   `flag` 패키지를 사용해 실행 시점에 `config.toml` 경로를 주입받음.
    *   컴파일 없이 환경(Local, Prod)을 교체할 수 있는 유연성 확보.
*   **Fail Fast 전략**:
    *   `cmd/app.go` 초기화 과정에서 의존성 주입 실패 시 즉시 `panic`을 발생시켜, 잘못된 상태로 서버가 켜지는 것을 방지.

### 3.2 Web Framework 도입 (Gin)
*   **Gin Framework 선택**:
    *   프로젝트의 메인 규칙은 Echo를 사용하지만, 본 `grpc-server` 모듈에서는 학습 및 gRPC Gateway 역할 수행 등을 위해 **Gin Framework**를 도입했습니다.
    *   `network/router.go`에서 `gin.Engine`을 초기화하고, `StartServer()` 메서드를 통해 HTTP 서버를 구동합니다.
*   **구현 내용**:
    *   `gin.New()`를 통한 엔진 초기화.
    *   `:8080` 포트로 서버 바인딩.
