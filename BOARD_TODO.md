# 📋 Board Module Learning Roadmap (CRUD)

이 문서는 Java Spring 개발자가 Go 언어로 간단한 게시판 CRUD를 구현하며 'Go 스러운 사고방식'을 익히기 위한 단계별 가이드입니다.

## 🎯 학습 목표
Spring의 `@Controller`, `@Service`, `@Repository`, `@Entity` 패턴을 Go의 구조체(Struct)와 인터페이스(Interface)로 직접 구현해보며, Magic(자동 설정)이 아닌 Explicit(명시적 코드)의 매력을 느껴봅니다.

---

## ✅ Step 1: 도메인 정의 (Domain Layer)
가장 먼저 데이터의 형태를 정의합니다.

- [ ] `internal/domain/post.go` 작성
    - [ ] `Post` 구조체 정의 (ID, Title, Content, CreatedAt)
    - [ ] JSON 태그 설정 (`json:"title"`)
    - [ ] **학습 포인트**:
        - Class 대신 Struct 사용
        - Lombok(@Getter, @Setter) 없이 필드에 직접 접근 (Public Field: 대문자로 시작)
        - Struct Tag를 이용한 JSON 직렬화/역직렬화 제어

## ✅ Step 2: 저장소 계층 구현 (Repository Layer)
데이터베이스 연동 전, 메모리(Map)를 활용하여 인터페이스를 먼저 정의하고 구현체를 만듭니다.

- [ ] `internal/repository/post_repository.go` 작성
    - [ ] `PostRepository` 인터페이스 정의 (Save, FindByID, FindAll, Update, Delete)
    - [ ] `memoryPostRepository` 구조체 구현 (Thread-safe를 위해 `sync.RWMutex` 활용)
    - [ ] **학습 포인트**:
        - **Interface**: Go의 인터페이스는 암시적(Implicit)입니다. (`implements` 키워드 없음)
        - **Pointer Receiver**: 상태를 변경하는 메서드에서 포인터(`*`) 사용법 (`func (r *Repo) Save...`)
        - **Concurrency**: `defer`와 `RWMutex`를 이용한 Go 스타일의 동시성 제어 기초

## ✅ Step 3: 서비스 계층 구현 (Service Layer)
핵심 비즈니스 로직을 담당합니다.

- [ ] `internal/service/post_service.go` 작성
    - [ ] `PostService` 구조체 정의 (Repository 인터페이스 주입)
    - [ ] `NewPostService` 생성자 함수 작성 (DI)
    - [ ] **학습 포인트**:
        - **Dependency Injection**: `@Autowired` 없이 생성자 함수로 의존성을 직접 주입하는 방법
        - **Error Handling**: `try-catch` 대신 `if err != nil`을 통해 에러를 값(Value)으로 처리하는 철학

## ✅ Step 4: 핸들러 계층 구현 (Handler Layer)
HTTP 요청을 처리하고 응답을 반환합니다. (Echo Framework 활용)

- [ ] `internal/handler/post_handler.go` 작성
    - [ ] `PostHandler` 구조체 정의 (Service 주입)
    - [ ] CRUD 메서드 구현 (`Create`, `Get`, `List`, `Update`, `Delete`)
    - [ ] Echo Context (`c echo.Context`) 활용
    - [ ] **학습 포인트**:
        - **Binding**: HTTP 요청 본문(JSON)을 구조체로 바인딩 (`c.Bind`)
        - **Routing**: 메서드를 일급 객체(First-class citizen)로 취급하여 라우터에 등록

## ✅ Step 5: 의존성 주입 및 서버 실행 (Wiring)
`main.go`에서 분리된 계층들을 조립합니다.

- [ ] `cmd/app/main.go` 수정
    - [ ] Repository -> Service -> Handler 순서로 의존성 객체 생성 및 주입
    - [ ] Echo 라우팅 그룹 설정 (`/posts`)
    - [ ] **학습 포인트**:
        - **Composition Root**: Spring Container가 해주던 일을 `main` 함수에서 직접 수행하며 애플리케이션의 전체 구조를 명확히 파악

---

## 🚀 Next Step (Challenge)
- [ ] 메모리 저장소를 실제 DB(PostgreSQL/MySQL)로 교체 (코드 수정 없이 구현체만 교체하는 인터페이스의 강력함 체험)
- [ ] DTO(Request/Response) 분리하여 엔티티 오염 방지
