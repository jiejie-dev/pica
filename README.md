# pica

> Pica is a restful automated testing tool and document generate tool written in golang.

## Install

```console
go get github.com/jeremaihloo/pica/cmd/pica
```

## Usage

```console
$ pica --help

usage: pica [<flags>] [<filename>] [<apiNames>...]

A command line for api test and doc generate

Flags:
  --help             Show context-sensitive help (also try --help-long and --help-man).
  --delay=DELAY      Delay after one api request.
  --output=OUTPUT    Output file.
  --filetype="pica"  The type of api file.
  --debug            Debug mode.
  --run              Run file.
  --convert          Convert file.
  --doc              Generate document for a api file.
  --server           Run as a document server
  --parse            Parse api file.
  --format           Format api file.

Args:
  [<filename>]  Api file.
  [<apiNames>]  Api names to excute

```

## TODO

- Document Generate
- Mock Server
- Api Document Server
- Api Document Version Control

## LICENSE

The MIT License (MIT)

Copyright (c) 2018 jeremaihloo