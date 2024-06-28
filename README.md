# HUP acceptable io.Writer

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/hupwriter)](https://pkg.go.dev/github.com/koron/hupwriter)
[![Actions/Go](https://github.com/koron/hupwriter/workflows/Go/badge.svg)](https://github.com/koron/hupwriter/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/hupwriter)](https://goreportcard.com/report/github.com/koron/hupwriter)

## Getting Started

### Install and Upgrade

```console
$ go get github.com/koron/hupwriter@latest
```

### Usage

```go
import (
    "github.com/koron/hupwriter"
    "log"
)

h := hupwriter.New("/var/log/myapp.log", "/var/pid/myapp.pid")
log.SetOutput(h)
```
