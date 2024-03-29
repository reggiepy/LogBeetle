# LogBeetle

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/reggiepy/LogBeetle?style=flat&color=success)
![GitHub Tag](https://img.shields.io/github/v/tag/reggiepy/LogBeetle?style=flat&color=success)
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

generate swagger UI
```bash
swag init -g cmd/log-writer/main.go
```

run
```bash
go run github.com/reggiepy/LogBeetle/cmd/LogBeetle server
```

build
```bash
go build github.com/reggiepy/LogBeetle/cmd/LogBeetle
```
## Architecture

```mermaid
sequenceDiagram
    participant Client
    participant HTTP
    participant NSQ
    participant LogBeetle
    Client-->>HTTP: Client send request to send log as a message
    loop HealthCheck
        LogBeetle->>LogBeetle: Consumers lookup nsq topic message
    end
    HTTP-->>NSQ: Send log to nsq topic
    Note right of LogBeetle: Consumer to handle log by nsq topic
    NSQ-->>LogBeetle: LogBeetle Consumer write log to file

```

