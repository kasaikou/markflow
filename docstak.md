# Scripts for `github.com/kasaikou/docstak` developpers

## download

Download dependencies

```yaml:docstak.yml
requires:
  exist: ["go.mod", "go.sum"]
```

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
go mod tidy && git diff --no-patch --exit-code go.sum &&
docstak download &&
gofmt -l . &&
docstak test
```
