# 🎈Java Spring Developer's `Go Playground`

> 이 저장소는 Java/Spring 개발자가 Go 언어의 생태계에 적응하고,
> 
> Spring에서의 개발 경험을 Go 환경에 어떻게 'Go 스럽게' 녹여낼 수 있는지 탐구하는 실습 프로젝트입니다.

### 🏗️ Project Architecture: Layered Approach
본 프로젝트는 Go 커뮤니티에서 가장 널리 권장되는 **Standard Project Layout**[📙](./ARCHITECTURE.md)을 기반으로, Spring의 계층화 구조를 재해석하여 설계했습니다.

### 폴더 구조 및 Spring 대응
```text
(예/board, ecommerce)/
├── cmd/                # 서비스의 진입점 (컴파일 대상)
│   └── app/main.go     # 애플리케이션 메인 실행 파일 및 진입점 (Spring Boot Main Class)
├── internal/           # 외부에서 참조 불가능한 내부 핵심 로직 (**강력한 캡슐화**)
│   ├── handler/        # Controller 역할 (HTTP/gRPC 입출력)
│   ├── service/        # Business Logic (핵심 도메인 로직)
│   ├── repository/     # Data Access 데이터 접근 로직 (DB, Cache)
│   └── domain/         # 핵심 엔티티 및 인터페이스 (Entity/Domain Model)
├── pkg/                # 외부 프로젝트에서도 재사용 가능한 유틸리티
└── configs/            # 환경 설정 및 설정 파일
```


### 왜 이 구조(Standard Layout)를 선택했는가?
실무에서는 생산성을 위해 프레임워크 지향적인 간결한 구조(예: [echo-boilerplate](https://github.com/alexferl/echo-boilerplate))를 사용하기도 하지만, 
본 프로젝트에서는 **Go의 인터페이스 기반 설계와 명시적 의존성 주입(DI) 원리**를 깊이 있게 파악하기 위해 표준 레이어드 아키텍처를 채택했습니다. 
이는 'Go 스러움'의 핵심인 **명확성**과 **테스트 용이성**을 확보하는 연습이기도 합니다.
[(이 Layout은 Go 스러운가?🙄)](./GO_IDIOMATIC_LAYOUT.md)

---

### 🛠️ Go Workspaces (go.work) 활용
Java의 Maven/Gradle 멀티 모듈 구조와 유사한 경험을 위해 go.work를 도입했습니다.

#### Go Workspaces도입 이유
1.  **멀티 모듈 관리**: 하나의 저장소 내에서 독립된 여러 모듈(`board`, `ecommerce` 등)을 동시에 개발하기 편리합니다.
2.  **로컬 의존성 해결**: 모듈 간의 의존성이 있을 때 `go.mod` 파일에 일일이 `replace` 구문을 추가하지 않고도 로컬에서 유연하게 참조할 수 있습니다.
3.  **Spring 멀티 프로젝트와 유사한 경험**: Java/Spring의 멀티 모듈 프로젝트(Maven/Gradle) 구조와 유사한 관리 환경을 제공하여 익숙한 방식으로 프로젝트를 구성할 수 있습니다.

### ✍️ Spring 개발자를 위한 Go 핵심 가이드 (Read Me First)

| Java / Spring | Go | 비고 |
| --- | --- | --- |
| **Annotation-based** | **Explicit Code** | 마법 같은 자동 설정 대신 명시적인 코드를 지향합니다. |
| **Exceptions** | **Explicit Error Handling** | `try-catch` 대신 `if err != nil`로 에러를 값으로서 처리합니다. |
| **Runtime Polymorphism** | **Static Duck Typing** | 인터페이스 구현을 명시하지 않아도 메서드 세트만 맞으면 성립합니다. |
| **Maven/Gradle** | **Go Modules & Workspaces** | 의존성 관리와 멀티 모듈 관리가 더 가볍고 빠릅니다. |

### 💡 Spring 개발자가 Go에서 가장 자주 하는 실수
Go에서는 interface를 구현하는 쪽에 두지 않고, 사용하는 쪽에 두는 것이 더 권장됩니다. 예를 들어 service가 repository를 필요로 한다면, service 패키지 안에 인터페이스를 정의하는 것이 가장 'Go 스러운' 추상화 방식입니다.

---

## 🛡️ Go Service Pattern 에러 헨들링 아키텍처

본 프로젝트는 비즈니스 로직의 독립성과 유지보수성을 위해 **에러 처리 패턴을 체계화**했습니다.

### 1. 핵심 목표 (Objective)
*   **관심사 분리 (Separation of Concerns):** Service Layer(비즈니스 로직)가 HTTP/Echo와 같은 전송 계층의 상세 구현(Status Code 등)을 알지 못하도록 격리합니다.
*   **도메인 독립성:** `board`, `ecommerce` 등 각 도메인 내부에서 에러를 정의하여, 향후 마이크로서비스 분리 시 독립성을 보장합니다.

### 2. 주요 구현 내용 (Implementation)

#### A. 제네릭 도메인 에러 정의 (`internal/domain/errors.go`)
프로토콜(HTTP, gRPC)에 종속되지 않는 추상화된 에러를 정의합니다.
*   **Standard Errors:** `ErrBadRequest`, `ErrNotFound`, `ErrInternalFailure` 등
*   **AppError 구조체:**
    *   `ServiceError`: 추상화된 에러 카테고리 (로직 분기용)
    *   `Detail`: 실제 발생한 기술적 에러 (서버 로깅용)
    *   `Message`: 클라이언트에게 전달할 안전한 메시지 (사용자 경험)

#### B. Service Layer: 에러 래핑 (Wrapping)
서비스 메서드는 에러 발생 시 `AppError`로 래핑하여 반환합니다. 이를 통해 상위 계층에 '에러의 종류'와 '사용자 메시지'를 명확히 전달합니다.

#### C. Handler Layer: 에러 번역 (Translation)
Echo 핸들러는 `errors.As`와 `errors.Is`를 사용하여 Service에서 넘어온 에러를 적절한 HTTP 응답으로 변환합니다.
*   `ErrBadRequest` -> HTTP 400
*   `ErrNotFound` -> HTTP 404
*   `ErrInternalFailure` -> HTTP 500

### 3. 도입 효과 (Benefits)
1.  **유연성:** 향후 gRPC 등 새로운 프로토콜 도입 시 Service 로직 수정 없이 Handler만 추가하면 됩니다.
2.  **보안성:** DB 에러 등 민감한 정보를 클라이언트에 노출하지 않고, `Message` 필드를 통해 안전한 메시지만 전달합니다.
3.  **유지보수성:** 에러 정의, 발생, 처리의 책임이 명확히 분리되어 디버깅과 확장이 용이합니다.

---

## 프로젝트 구성

- **board**: Echo 프레임워크 기반의 계층형 아키텍처(Layered Architecture) 게시판 프로젝트
    - [👉 Board 모듈 개발 로드맵 (CRUD 따라하기)](./BOARD_TODO.md)
- **ecommerce**: 향후 추가 예정인 이커머스 프로젝트

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
