# HUP acceptable io.Writer

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/hupwriter)](https://pkg.go.dev/github.com/koron/hupwriter)
[![Actions/Go](https://github.com/koron/hupwriter/workflows/Go/badge.svg)](https://github.com/koron/hupwriter/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/hupwriter)](https://goreportcard.com/report/github.com/koron/hupwriter)

Package hupwriter provides wrapper of `os.File`.

The wrapper will close and reopen the file when a HUP signal is received.
It allows easier log file rotation.

By logging to `hupwriter.HupWriter` you can create log files that can be used
with log rotation management programs such as logrotate, newsyslog, or so.

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

// TODO: log something.
```

```console
## Rename the log file to rotate
# mv /var/log/myapp.log /var/log/myapp.0

## Send HUP signal to the process. It closes and reopens myapp.log file.
# kill -HUP $(cat /var/pid/myapp.pid)

## (OPTIONAL) Compress rotated log file
# gzip /var/log/myapp.0
```
