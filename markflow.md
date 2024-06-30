# Scripts for `github.com/kasaikou/markflow` developpers

## hello_world

Echo "Hello World"

```sh
echo "Hello World, docstak!"
```

## download

Download dependencies

```yaml:docstak.yml
requires:
  file:
    exist: ["go.mod", "go.sum"]
skips:
  file:
    not-changed: ["go.sum"]
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

## ci

Running on GitHub Actions, local, and so on.

```yaml:docstak.yml
previous: [ci/fmt, ci/depends, ci/coverage-test]
```

### ci/depends

```yaml:docstak.yml
skips:
  file:
    not-changed: ["**.go", "go.sum", "go.mod"]
```

```sh
go mod tidy &&
git diff --no-patch --exit-code go.sum
```

### ci/fmt

```sh
gofmt -l .
```

### ci/coverage-test

```yaml:docstak.yml
previous: [ci/coverage-test/go, download]
```

#### ci/coverage-test/go

```yaml:docstak.yml
skips:
  file:
    not-changed: ["**.go", "go.sum", "go.mod"]
```

```sh
go test -coverprofile=coverage.txt -covermode=atomic ./...
```
