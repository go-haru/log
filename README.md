# log

[![Go Reference](https://pkg.go.dev/badge/github.com/go-haru/log.svg)](https://pkg.go.dev/github.com/go-haru/log)
[![License](https://img.shields.io/github/license/go-haru/log)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/go-haru/log.svg?style=flat-square)](https://github.com/go-haru/log/releases)
[![Go Test](https://github.com/go-haru/log/actions/workflows/go.yml/badge.svg)](https://github.com/go-haru/log/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-haru/log)](https://goreportcard.com/report/github.com/go-haru/log)

Here's a simple and extensible logging library for Go.

## Features
- Simple API for basic logging
- Data field collecting and json marshaling
- Replaceable underlying implementation

## Interface

### Level

`type Level uint`

`const (`

 - `DebugLevel Level = iota`
 - `InfoLevel`
 - `WarningLevel`
 - `ErrorLevel`
 - `FatalLevel`

`)`

### Logger

`type Logger interface {`

- `Debug(v ...any)`
- `Info(v ...any)`
- `Warn(v ...any)`
- `Error(v ...any)`
- `Fatal(v ...any)` // will cause os.Exit(-1)
- `Panic(v ...any)` // will throw panic
- `Print(v ...any)` // use level from `WithLevel`
- `Debugf(format string, v ...any)`
- `Infof(format string, v ...any)`
- `Warnf(format string, v ...any)`
- `Errorf(format string, v ...any)`
- `Fatalf(format string, v ...any)` // will cause os.Exit(-1)
- `Panicf(format string, v ...any)` // will throw panic
- `Printf(format string, v ...any)` // use level from `WithLevel`
- `With(v ...field.Field) Logger`
  get cloned Logger with data field
- `WithName(name string) Logger`
  get cloned Logger with new name prefix
- `WithLevel(level Level) Logger`
  get cloned Logger with new level
- `AddDepth(depth int) Logger`
  get cloned Logger with increased caller skip count
- `Standard() *log.Logger`
  create a native logger, level comes from `WithLevel`
- `Flush() error`
  force sync buffer now

`}`

## Function
 
- `func Debug(m ...any)`
- `func Info(m ...any)`
- `func Warn(m ...any)`
- `func Error(m ...any)`
- `func Panic(m ...any)`
- `func Fatal(m ...any)`
- `func Debugf(format string, m ...any)`
- `func Infof(format string, m ...any)`
- `func Warnf(format string, m ...any)`
- `func Errorf(format string, m ...any)`
- `func Panicf(format string, m ...any)`
- `func Fatalf(format string, m ...any)`
- `func Printf(format string, m ...any)`
- `func With(v ...field.Field) Logger` // same as Logger.With
- `func WithName(name string) Logger`  // same as Logger.WithName
- `func Use(l Logger)`                 // replace global logger ( not thread safe )
- `func Current() Logger`              // clone global logger


## Usage

### Basic example

```go
package main

import (
  "encoding/json"
  "os"

  "github.com/go-haru/field"
  "github.com/go-haru/log"
)

package func main() {
  // replace logger ( optional, we have inherent implementation )
  log.With( /* your own implementation */)

  // normal text
  log.Info("starting up")

  // export logger
  var logger = log.Current()

  var configData struct{}
  var configFile = "./config.json"

  // add data field
  logger = logger.With(field.String("configFile", configFile))

  // mixed
  if h, err := os.Open(configFile); err != nil {
    logger.With(field.Error("cause", err)).Error("cant open file")
  } else if err = json.NewDecoder(h).Decode(&configData); err != nil {
    // fatal, will cause os.exit(-1)
    logger.With(field.Error("cause", err)).Fatal("cant parse file")
  }
}

```

### Context injection example

as an example, we can make session logger for http request with tracing id:

```go
package main

import (
  "context"
  "net"
  "net/http"

  "github.com/google/uuid"
  "github.com/go-haru/field"

  "github.com/go-haru/log"
)

func main() {
  (&http.Server{
    BaseContext: func(net.Listener) context.Context {
      // make logger
      var logger = log.With(field.Stringer("traceId", uuid.New()))
      // inject context
      return log.Context(context.Background(), logger)
    },
    Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
      // fetch logger
      var logger = log.C(request.Context())

      logger.Debugf("url: %v", request.URL.String())
    }),
  }).ListenAndServe()
}

```

## Contributing

For convenience of PM, please commit all issue to [Document Repo](https://github.com/go-haru/go-haru/issues).

## License

This project is licensed under the `Apache License Version 2.0`.

Use and contributions signify your agreement to honor the terms of this [LICENSE](./LICENSE).

Commercial support or licensing is conditionally available through organization email.
