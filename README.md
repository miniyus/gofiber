# miniyus: Go-fiber Template

## go lang with fiber framework

Go 언어의 Fiber 웹프레임워크를 활용하여 구조화한 저장소입니다.

- gorm 사용
    - 현재 리포지토리에서는 postgreSQL만 사용중

## Install

```shell
git clone https://github.com/miniyus/gofiber.git

# 도커 컨테이너 및 gorm DB 드라이버는 postgreSQL을 사용 중입니다.
docker compose up -d --build 

# 로컬 실행 시
make start
# or
make build
./build/gofiber
```

## Dot Env

```shell
APP_ENV=develop
APP_NAME=go-fiber
APP_PORT=8080

TIME_ZONE=Asia/Seoul

GO_GROUP=1000
GO_USER=1000

DB_HOST=go-pgsql
DB_DATABASE=go_fiber
DB_PORT=5432
DB_USERNAME=?
DB_PASSWORD=?
DB_AUTO_MIGRATE=true

```

## Directory Structure

```shell
/project-root
|-- build: build된 파일이 저장
|-- cmd: go 실행 및 makefile을 통한 cli 명령을 실행할 수 있는 main.go 파일
|-- config: 설정
|-- data: local data 저장소
|-- docker: docker 컨테이너 관련
|-- internal: api를 실질적으로 구현하는 곳 입니다.
  |-- app: api 요청 시, 수행되는 코드들
    |-- dto: DTO 정의 및 매핑 함수 정의
    |-- factory: Handler 생성을 위한 팩토리 패턴 적용
    |-- handler: 요청을 받고 응답을 해준다.
    |-- service: handler 요청의 비즈니스 로직 처리
    |-- repositroy: db, entity를 통해 데이터 CRUD 동작 수행
    |-- routes: 그룹화된 API 적용을 위한 서브 라우터
  |-- core: 공통 기능
   |-- api_error: api error, error response 관련 기능
   |-- auth: 인증 관련 기능
   |-- container: go fiber앱을 감싸고 필요한 미들웨어 및 라우팅 등을 컨테이너를 통해서 제어 할 수 있게
   |-- context: fiber context Locals를 통해 가져올 항목들을 미리 정의
   |-- database: database, gorm 연결
   |-- logger: 로거
   |-- register: 컨테이너에 필요한 미들웨어, 구조체, 함수등을 등록
  |-- entity: db 스키마를 가진 구조체 집합
  |-- routes: 라우팅
  |-- utils: 유틸 함수들
|-- pkg: 독립적인 기능을 수핼 할 수 있는 기능들의 집합입니다.
 
```

### config

- 설정 관리
- go 언어의 구조체를 활용하여 관리

```go
import (
"github.com/miniyus/go-fiber/config"
)

config.GetConfigs()
```
### internal

### core.container

```go
// fiber 앱과 Configs 구조체를 인수로 받는다.
func NewContainer(app *fiber.App, config *config.Configs)


```

**Register 패키지**
> [internal/core/register/register.go](internal/core/register/register.go)

총 4개의 함수로 구성

- boot(): 우선순위가 높은 핸들러 혹은 미들웨어 등록
- middlewares(): 일반적인 글로벌 미들웨어 등록
- routes(): 라우트 세팅
- Register(): 위 private 함수들을 순서대로 실행 시켜주는 public 함수

