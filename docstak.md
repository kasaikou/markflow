# Scripts for `github.com/kasaikou/docstak` developpers

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
