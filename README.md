# SkyQuality Meter

[![](https://img.shields.io/github/actions/workflow/status/gaker/unihedron-sqm/go.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/gaker/unihedron-sqm/actions?query=workflow%3ATest) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gaker_unihedron-sqm&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gaker_unihedron-sqm) [![Go Reference](https://pkg.go.dev/badge/github.com/gaker/unihedron-sqm.svg)](https://pkg.go.dev/github.com/gaker/unihedron-sqm) [![Go Report Card](https://goreportcard.com/badge/github.com/gaker/unihedron-sqm)](https://goreportcard.com/report/github.com/gaker/unihedron-sqm)

Go implementation for the [Unihedron Sky Quality meter](http://www.unihedron.com/).

## Installation

```
go get -u github.com/gaker/unihedron-sqm
```

## Basic Usage

```
package main

import (
    "github.com/gaker/unihedron-sqm"
)

func main() {
    // take a reading
    host := "192.168.1.22"
    port := "10001"

    monitor := sqm.New(&sqm.Config{
        Tcp: &sqm.Tcp{
            Host: &host,
            Port: &port,
        },
    })

    if err := monitor.Dial(); err != nil {
        log.Fatalf("error dialing: %s", err)
    }

    if err := monitor.Send(sqm.CmdAveragedReading); err != nil {
        log.Fatalf("error sending command: %s", err)
    }

    reading := sqm.Reading{}
    if err := monitor.Read(&reading); err != nil {
        log.Fatalf("error reading: %s", err)
    }

    fmt.Printf("reading: %v", reading.Reading)
}
```

## TODO:
* [ ] SQM-LU [USB/Serial]
* [ ] SQM-LR RS232