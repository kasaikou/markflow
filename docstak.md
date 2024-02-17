# Scripts for `github.com/kasaikou/docstak` developpers

## download

Download dependencies

```sh
docstak download:go
```

### download:go

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

## ci-lint-test

Running on GitHub Actions, local, and so on.

```sh
docstak download &&
gofmt -l . &&
docstak test
```
