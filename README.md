# log

Go logging library wrapped [zap](https://github.com/uber-go/zap)

[![Build Status](https://travis-ci.org/FeiniuBus/log.svg?branch=master)](https://travis-ci.org/FeiniuBus/log)

## Installation

`go get -u github.com/FeiniuBus/log`

## Normal logger

```Go
package main

import (
    "fmt"

    "github.com/FeiniuBus/log"
)

func main() {
    logger, err := log.New(false)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer logger.Sync()

    logger.With("url", "http://www.baidu.com").Warn("failed to fetch URL")
}
```

## Logstash

```Go
package main

import (
    "fmt"

    "github.com/FeiniuBus/log"
)

func main() {
    logger, err := log.NewLogstash(false, "host", port)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer logger.Sync()

    logger.With("url", "http://www.baidu.com").Warn("failed to fetch URL")
}
```