# LogBeetle

[![go version](https://img.shields.io/github/go-mod/go-version/reggiepy/LogBeetle?color=success&filename=go.mod&style=flat)](https://github.com/reggiepy/win_server_gen)
[![release](https://img.shields.io/github/v/tag/reggiepy/win_server_gen?color=success&label=release)](https://github.com/reggiepy/win_server_gen)
[![build status](https://img.shields.io/badge/build-pass-success.svg?style=flat)](https://github.com/reggiepy/win_server_gen)
[![License](https://img.shields.io/badge/license-GNU%203.0-success.svg?style=flat)](https://github.com/reggiepy/win_server_gen)
[![Go Report Card](https://goreportcard.com/badge/github.com/reggiepy/win_server_gen)](https://goreportcard.com/report/github.com/reggiepy/win_server_gen)

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

## Architecture
