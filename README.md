# pica

> Pica is a restful automated testing tool and document generate tool written in golang.

It's inspired deeply by [frank](https://github.com/txthinking/frank).

## Features

- Base api test (POST, GET, PUT, DELETE, PATCH)
- Generate api document to markdown file.
- Benchmark webapi.(TODO)
- Serve api document as a website.(TODO)
    - Custom theme or css for this website(TODO).
    - Api version controls, automated version release note.(TODO)
    - Api version diff to show.(TODO)

## Status

It's under development.

## Install

```console
go get github.com/jeremaihloo/pica/cmd/pica
```

## Usage

```console
$ pica --help

usage: pica [<flags>] <command> [<args> ...]

A command line for api test and doc generate.

Flags:
  --help   Show context-sensitive help (also try --help-long and --help-man).
  --debug  Debug mode.

Commands:
  help [<command>...]
    Show help.

  run [<flags>] [<filename>] [<apiNames>...]
    Run api file.

  format [<flags>] [<filename>]
    Format api file.

  serve [<flags>]
    Run a document website.

  init [<filename>] [<template>]
    Init a new api file from template.

  config [<flags>]
    Config pica.


```

## TODO

- Document Generate
- Api Document Server
- Api Document Version Control

## LICENSE

The MIT License (MIT)

Copyright (c) 2018 jeremaihloo