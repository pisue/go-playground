## 🛠 Prerequisites for Go gRPC Development (Windows)

이 프로젝트는 Go 언어와 gRPC를 사용하며, Windows 환경에서 원활한 코드 생성을 위해 아래 도구들의 설치 및 환경 변수 설정이 필수적입니다.

### 1. Protocol Buffer 컴파일러 (protoc) 설치

`protoc`는 `.proto` 파일을 읽어 코드를 생성하는 메인 엔진입니다.

1. [Protobuf GitHub Releases](https://github.com/protocolbuffers/protobuf/releases)에서 `protoc-xx.x-win64.zip` 파일을 다운로드합니다.
2. `C:\protoc` 폴더를 생성하고 압축을 풉니다. (가급적 한글이 포함되지 않은 경로를 권장합니다.)
* 최종 경로 예시: `C:\protoc\bin\protoc.exe`


3. **환경 변수 설정:**
* 시스템 환경 변수 편집 -> `Path` 변수 선택 후 [편집] -> [새로 만들기]
* `C:\protoc\bin` 경로를 추가하고 **맨 위로 이동**시킵니다.



### 2. Go gRPC 플러그인 설치

`protoc`가 Go 코드를 생성할 수 있도록 돕는 플러그인들을 설치합니다.

```powershell
# Go 프로토콜 버퍼 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Go gRPC 플러그인 설치
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

### 3. 환경 변수(PATH) 최종 확인

Windows의 경우, 위 플러그인들은 보통 `%USERPROFILE%\go\bin` 폴더에 설치됩니다. 명령어를 어디서든 실행할 수 있도록 아래 두 경로가 모두 **시스템 환경 변수 `Path`**에 등록되어 있어야 합니다.

| 순위 | 등록할 경로 예시 | 비고 |
| --- | --- | --- |
| 1 | `C:\protoc\bin` | protoc 본체 실행을 위함 |
| 2 | `%USERPROFILE%\go\bin` | Go 플러그인 실행을 위함 |

> **⚠️ 중요:** 환경 변수를 설정한 후에는 **반드시 터미널(PowerShell)과 IDE(GoLand, VS Code 등)를 완전히 종료한 후 다시 실행**해야 변경 사항이 반영됩니다.

---

## 🚀 Code Generation

설치가 완료되었다면, 프로젝트 루트 디렉토리에서 아래 명령어를 실행하여 gRPC 코드를 생성합니다.

```powershell
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       [proto_파일_경로]

```

*(PowerShell에서 줄바꿈 입력 시 ``` 문자를 사용하거나, 모든 명령어를 한 줄로 이어 붙여 실행하세요.)*

---

### 💡 Troubleshooting

* **'protoc'은(는) 내부 또는 외부 명령...이 아닙니다**: `C:\protoc\bin` 경로가 환경 변수에 정확히 등록되었는지, 그리고 터미널을 재시작했는지 확인하세요.
* **GoLand 터미널에서만 안 되는 경우**: IDE가 이전 환경 변수 세션을 유지하고 있을 확률이 높습니다. GoLand를 완전히 종료 후 다시 실행하거나, 새 터미널 탭을 열어보세요.

---
맥북(macOS)이나 리눅스 환경에서 작업할 동료나 미래의 화평님을 위해, **UNIX 계열 환경에서의 설정법**을 README 하단에 추가해 두면 완벽합니다. 윈도우와는 설치 도구와 환경 변수 설정 파일이 다르다는 점이 핵심입니다.

아래 내용을 기존 README 내용 아래에 이어서 붙여넣으세요.

---

## 🍎 Prerequisites for gRPC Development (macOS / Linux)

UNIX 계열 환경에서는 패키지 매니저를 통해 훨씬 간편하게 도구를 설치할 수 있습니다.

### 1. Protocol Buffer 컴파일러 (protoc) 설치

맥북에서는 **Homebrew**를 사용하는 것이 가장 권장됩니다.

```bash
# Homebrew를 통한 설치
brew install protobuf

# 설치 확인
protoc --version

```

### 2. Go gRPC 플러그인 설치 (Windows와 동일)

프로토콜 버퍼 파일을 Go 코드로 변환해 주는 플러그인들을 설치합니다.

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

### 3. 환경 변수(PATH) 설정 (zsh 기준)

맥북의 기본 쉘인 `zsh`에서 `go/bin` 경로를 인식할 수 있도록 설정 파일(~/.zshrc)에 추가해야 합니다.

1. 터미널에서 설정 파일을 엽니다.
```bash
nano ~/.zshrc

```


2. 파일 하단에 아래 내용을 복사해서 붙여넣습니다.
```bash
# Go Bin Path
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

```


3. `Ctrl + O` (저장), `Enter`, `Ctrl + X` (종료)를 순서대로 누릅니다.
4. 변경 사항을 즉시 적용합니다.
```bash
source ~/.zshrc

```



---

## 🚀 Code Generation (macOS / Linux)

맥북 터미널(bash/zsh)에서는 줄바꿈 기호로 역슬래시(`\`)를 사용합니다.

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       [proto_파일_경로]

```

---

### 💡 Tips for Mac Users

* **M1/M2/M3 (Apple Silicon) 사용자**: Homebrew로 설치한 `protobuf`가 정상적으로 동작하지 않는다면, `brew doctor`를 통해 경로 설정을 확인해 보세요.
* **권한 문제**: `go install` 시 권한 에러가 발생한다면 `sudo`를 사용하지 말고, 위 환경 변수(`GOPATH`) 설정이 정확히 되어 있는지 먼저 확인하는 것이 좋습니다.

---