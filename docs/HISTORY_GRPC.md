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

### 3.3 gRPC 환경 구축 및 초기 서비스 정의
*   **프로토콜 버퍼(Protocol Buffers) 설정**:
    *   Windows 및 macOS/Linux 환경을 위한 `protoc` 설치 및 Go 플러그인(`protoc-gen-go`, `protoc-gen-go-grpc`) 설정 가이드(`docs/gRPC_SETTING.md`) 작성.
*   **AuthService 정의 (`auth.proto`)**:
    *   인증 관련 RPC 메서드 `CreateAuth`, `VerifyAuth` 설계.
    *   토큰 생성을 위한 데이터 구조 `AuthData` 및 응답 메시지 정의.
*   **코드 생성**:
    *   `protoc` 명령어를 통해 `.proto` 파일로부터 Go 소스 코드(`auth.pb.go`, `auth_grpc.pb.go`) 생성 성공.
*   **의존성 추가**:
    *   `grpc`, `protobuf`, `paseto` (토큰 관리), `gin` 등 핵심 라이브러리 의존성 주입.

### 3.4 gRPC 클라이언트 설계 및 상태 공유 개념 학습 (State Management)
gRPC 서비스 호출을 담당하는 클라이언트(`client.go`)를 구현하며, **구조체 필드와 메서드 인자의 역할 분리**에 대해 깊이 있게 학습했습니다.

#### A. 구조체와 설정 공유 (Context vs. State)
초기 의문점: *"Config 포인터를 공유하면 모든 요청이 똑같은 값을 바라보게 되어 데이터가 섞이지 않을까?"*

*   **해결 및 정립된 개념**:
    *   **도구(Tools/Resources)**: `Config`, `ClientConn`(연결 객체), `PasetoMaker`(토큰 생성기) 등은 **불변(Immutable)**하거나 **공유되어야 하는 자원**입니다. 이는 `GRPCClient` 구조체의 필드로 관리하며, 포인터로 공유하여 메모리 효율과 일관성을 유지합니다.
    *   **재료(Data/Request)**: 사용자별 요청 데이터(`Address`, `Token` 등)는 **가변(Mutable)**하며 요청마다 다릅니다. 이는 구조체 필드가 아닌 **메서드의 인자(Parameter)**로 전달되어 스택 메모리에서 처리되고 소멸합니다.

#### B. 구현 내용 (`GRPCClient`)
*   **구조체 정의**:
    ```go
    type GRPCClient struct {
        client      *grpc.ClientConn        // 연결 객체 (공유)
        authClient  auth.AuthServiceClient  // gRPC 스텁 (공유)
        pasetoMaker *paseto.PasetoMaker     // 토큰 생성기 (공유)
    }
    ```
*   **생성자 (`NewGRPCClient`)**:
    *   `grpc.Dial`을 통해 서버와 연결을 맺고 `insecure` 옵션으로 개발 환경 통신을 설정.
    *   `Config` 구조체를 확장하여 gRPC 서버 URL(`GRPC.URL`)을 외부 설정 파일에서 주입받도록 개선.
*   **Paseto 모듈 개선**:
    *   `NewPasetoMaker` 생성자의 인자를 값(`config.Config`)에서 포인터(`*config.Config`)로 변경하여 일관성 확보.
