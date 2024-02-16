# Scripts for `github.com/kasaikou/docstak` developpers

```yaml:docstak.yml
- shells
  - defaults: [`sh`, `powershell`]
- dotenvs
  - `.env` required
- envs
  - `.env` 
```

## download

Download dependencies

```sh
go mod download
```

## test

Run go test

```sh
DOCSTAK_TEST_WORKSPACE_DIR=$(pwd) go test ./...
```

## fmt

Format source codes

```sh
go fmt ./...
```
