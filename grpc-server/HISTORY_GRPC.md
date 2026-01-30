# gRPC Server 학습 및 개발 이력

이 문서는 `grpc-server` 프로젝트를 진행하며 학습한 내용과 개발 변경 사항을 기록합니다.
기존에 진행했던 **Echo 프레임워크 기반의 HTTP 웹 서버** 프로젝트와 달리, 본 프로젝트는 **gRPC 프로토콜의 핵심 원리와 동작 방식**을 깊이 있게 이해하는 것을 최우선 목표로 합니다.

## 1. 프로젝트 개요 및 목적

*   **목표**: gRPC, Protocol Buffers(Protobuf)의 이해 및 Go 언어에서의 gRPC 서버 구현 학습.
*   **차별점**: 프레임워크의 편리한 기능(Magic)에 의존하기보다, **Pure Go**에 가까운 방식으로 직접 의존성을 주입하고 서버를 구성하여 내부 동작 흐름을 명확히 파악합니다.

## 2. 기존 Echo 프로젝트와의 아키텍처 및 구현 비교

나중에 다시 보았을 때 기존 HTTP(REST) 방식과 gRPC 방식, 그리고 구조적 설계의 차이를 이해하기 위해 상세히 기록합니다.

### 2.1 통신 프로토콜 및 진입점 (Entry Point)

| 구분 | Echo 프로젝트 (HTTP/REST) | gRPC 프로젝트 (RPC) |
| :--- | :--- | :--- |
| **기반 프로토콜** | HTTP/1.1 (주로 사용), JSON 텍스트 포맷 | HTTP/2, Protocol Buffers (바이너리 포맷) |
| **API 정의** | 코드 내 라우팅 (`e.GET("/users", ...)`), Swagger 등 | `.proto` 파일 (IDL)로 명세 정의 후 코드 생성 |
| **핸들러 위치** | `handler` 패키지 (Controller 역할) | `network` 패키지 (gRPC Server Implementation) |
| **데이터 바인딩** | `c.Bind(&dto)` (런타임 리플렉션/JSON 파싱) | `protoc`로 생성된 Go 구조체 사용 (컴파일 타임 타입 안전) |

> **핵심 차이**: Echo는 URL 경로와 HTTP 메서드로 라우팅하지만, gRPC는 `.proto`에 정의된 **서비스 인터페이스의 메서드를 직접 호출**하는 것처럼 동작합니다.

### 2.2 아키텍처 구조 및 의존성 주입 (Dependency Injection)

기존 프로젝트가 프레임워크의 기능이나 외부 DI 라이브러리(예: `uber-go/fx`)를 사용해 "설정보다 관례" 혹은 "자동화"를 추구했다면, 이번 프로젝트는 **명시성(Explicitness)**을 극대화했습니다.

#### A. 의존성 주입 (Wiring) 방식
*   **기존 방식 (Implicit/Automated)**:
    *   DI 컨테이너가 자동으로 타입을 찾아 주입하거나, 전역 변수 등을 사용할 수도 있음.
    *   편리하지만, "누가 누구를 참조하는지" 코드만 보고 즉시 파악하기 어려울 때가 있음.
*   **현재 방식 (Explicit Manual Wiring - `cmd/app.go`)**:
    *   `NewApp` 함수 내에서 `if-else` 체이닝을 통해 객체 생성 순서를 강제함.
    *   **순서**: `Config` -> `Repository` -> `Service` -> `Network`
    *   코드 자체가 설계도 역할을 하여, 데이터 흐름과 의존 관계가 한눈에 보임.

#### B. 레이어 구성
```text
[요청 흐름]
Client -> (gRPC/Proto) -> Network Layer -> Service Layer -> Repository Layer -> DB
```

*   **Network Layer (`network/`)**:
    *   기존의 `handler`나 `controller`에 해당합니다.
    *   하지만 이곳은 gRPC가 생성해준 `UnimplementedServer` 인터페이스를 구현하는 곳입니다.
    *   HTTP Status Code 대신 gRPC Status Code를 다룹니다.
*   **Service Layer (`service/`)**:
    *   비즈니스 로직은 동일하게 유지됩니다.
    *   gRPC와 무관하게 순수 Go 로직으로 작성하여, 향후 프로토콜이 바뀌어도 재사용 가능하도록 격리합니다.

## 3. 개발 진행 로그

### 3.1 초기 스캐폴딩 (Layered Architecture)
*   **Flag 기반 설정 관리**:
    *   `flag` 패키지를 사용해 실행 시점에 `config.toml` 경로를 주입받음.
    *   컴파일 없이 환경(Local, Prod)을 교체할 수 있는 유연성 확보.
*   **Fail Fast 전략**:
    *   `cmd/app.go` 초기화 과정에서 의존성 주입 실패 시 즉시 `panic`을 발생시켜, 잘못된 상태로 서버가 켜지는 것을 방지.

---