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

## ci

Running on GitHub Actions, local, and so on.

```yaml:docstak.yml
previous: [ci/fmt, ci/depends, ci/test]
```

## ci/depends

```sh
go mod tidy && git diff --no-patch --exit-code go.sum
```

## ci/fmt

```sh
gofmt -l .
```

## ci/test

```yaml:docstak.yml
previous: [ci/test/go]
```

## ci/test/go

```sh
go test -coverprofile=coverage.txt -covermode=atomic ./...
```

## ci/test/send-coverage

Send to Codecov's coverage report.

```yaml:docstak.yml
previous: [ci/test/go]
```

```sh
codecov upload-process -f coverage.txt -t ${CODECOV_TOKEN}
```
