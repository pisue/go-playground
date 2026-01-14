# Standard Go Project Layout

이 프로젝트는 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)을 참고하여 구성되었습니다.
Go 생태계에서 널리 사용되는 디렉토리 구조 패턴이며, 프로젝트의 유지보수성과 확장성을 높이는 데 도움을 줍니다.

Java Spring 개발자에게 익숙한 계층형 아키텍처와 유사하게 매핑할 수 있습니다.

## 디렉토리 구조 설명

### `/cmd`
애플리케이션의 진입점(Main)이 위치하는 곳입니다.
- 이 디렉토리 내의 파일은 보통 작아야 하며, `main` 패키지를 포함합니다.
- 비즈니스 로직을 여기에 두지 않고, `internal`이나 `pkg` 디렉토리의 코드를 호출하고 설정(DI, Configuration)하는 역할만 수행합니다.
- **Spring 비유**: `@SpringBootApplication`이 있는 메인 클래스와 유사합니다.
- 예: `cmd/app/main.go`

### `/internal`
외부(다른 저장소/모듈)에서 임포트(Import)할 수 없는 비공개 애플리케이션 및 라이브러리 코드입니다.
- Go 컴파일러가 강제하는 규칙으로, 이 디렉토리 내의 패키지는 프로젝트 내부에서만 접근 가능합니다.
- 프로젝트의 핵심 비즈니스 로직, 도메인 모델, 핸들러 등이 이곳에 위치합니다.

#### 하위 디렉토리 구조 예시
- **`handler`**: HTTP 요청을 처리하고 응답을 반환합니다. (**Spring**: `@Controller`, `@RestController`)
- **`service`**: 비즈니스 로직을 수행합니다. (**Spring**: `@Service`)
- **`repository`**: 데이터베이스 접근을 담당합니다. (**Spring**: `@Repository`, `Repository` Interface)
- **`model`**: 데이터 구조체(Entity, DTO)를 정의합니다. (**Spring**: `@Entity`, DTO POJO)

### `/pkg`
외부 프로젝트에서도 사용할 수 있는 공용 라이브러리 코드입니다.
- 다른 프로젝트에서도 재사용 가능한 유틸리티나 헬퍼 함수 등을 이곳에 둡니다.
- `internal`과 달리 외부에서 `import`가 가능하므로, 공개 라이브러리로서의 성격을 가집니다.
- **Spring 비유**: `common-libs`와 같은 공통 모듈이나 유틸리티 클래스(`StringUtils` 등)에 해당합니다.
- 예: `pkg/utils`

### `/configs`
설정 파일이나 설정 템플릿을 저장합니다.
- `yaml`, `json`, `toml`, `.env` 등의 설정 파일이 위치합니다.
- **Spring 비유**: `application.yml` 또는 `application.properties`가 위치하는 `src/main/resources`와 유사합니다.

### `/api` (Optional)
OpenAPI/Swagger 명세, 프로토콜 버퍼(Protobuf) 파일, JSON 스키마 등 API 정의 파일을 둡니다.

### `/web` (Optional)
정적 웹 자원(HTML, JS, CSS)이나 템플릿 파일 등을 둡니다.
- **Spring 비유**: `src/main/resources/static` 또는 `templates` 디렉토리.
