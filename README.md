# LogBeetle

[![go version](https://img.shields.io/github/go-mod/go-version/reggiepy/LogBeetle?color=success&filename=go.mod&style=flat)](https://github.com/reggiepy/LogBeetle)
[![release](https://img.shields.io/github/v/tag/reggiepy/LogBeetle?color=success&label=release)](https://github.com/reggiepy/LogBeetle)
[![build status](https://img.shields.io/badge/build-pass-success.svg?style=flat)](https://github.com/reggiepy/LogBeetle)
[![License](https://img.shields.io/badge/license-GNU%203.0-success.svg?style=flat)](https://github.com/reggiepy/LogBeetle)
[![Go Report Card](https://goreportcard.com/badge/github.com/reggiepy/LogBeetle)](https://goreportcard.com/report/github.com/reggiepy/LogBeetle)

## Installation

```bash
git clone https://github.com/reggiepy/LogBeetle.git
cd LogBeetle
go mod tidy
```

## Usage

```bash
go run cmd/log-writer/main.go
```

生成swagger UI
```bash
swag init -g cmd/log-writer/main.go
```

build cli
```bash
go run build github.com/reggiepy/LogBeetle/cmd/cli
```

编译cli
```bash
cd cmd\cli
go run .
go build .
```

编译cli
```bash
go run build github.com/reggiepy/LogBeetle/cmd/cli
```
## Architecture
