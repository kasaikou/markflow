# Scripts for `github.com/kasaikou/docstak` developpers

## hello_world

Echo "Hello World"

```sh
echo "Hello World, docstak!"
```

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

```yaml:docstak.yml
previous: [download]
```

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

```yaml:docstak.yml
previous: [download]
```

```sh
go mod tidy && git diff --no-patch --exit-code go.sum &&
gofmt -l . &&
docstak test
```
