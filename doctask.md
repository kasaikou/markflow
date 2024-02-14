# Scripts for `github.com/kasaikou/run-scripts` developpers

```yaml:run-scripts.yml
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
go test ...
```

## fmt

Format source codes

```sh:/bin/bash
go fmt ...
```
