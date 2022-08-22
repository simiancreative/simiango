# Simian Go

[![Coverage Status](https://coveralls.io/repos/github/simiancreative/simiango/badge.svg?branch=master)](https://coveralls.io/github/simiancreative/simiango?branch=master)
[![tests](https://github.com/simiancreative/simiango/workflows/CI/badge.svg)](https://github.com/simiancreative/simiango/actions)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/simiancreative/simiango) 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/simiancreative/simiango/blob/master/LICENSE)


A multi purpose tool set for golang applications. Tools include:

- Dotenv Config
- Logger
- Http Server
  - Private Routes
  - Streaming Response
  - Request Parameter Parsing
- Service Handlers
- Reversible and Non-Reversible encryption
- JWT Token Generation and Validation
- Struct Validation
- Various Service Implementations meant to abstract the usage into a simple interface
  - Mysql
  - Mssql
  - Postrges
  - Sql null types without valid checking 😅
  - Redis
  - AMQP
  - Kafka

## Setup

Create a go application and import simian go packages as required. See the
[sample application](https://github.com/simiancreative/simiango/blob/master/examples/main.go)
for usage.

## CLI

install
```
go install github.com/simiancreative/simiango/simian-go@latest
```

run
```
simian-go -h
```

### How to generate your crypt-keeper key

Generate a base64 encoded 32 bit key and copy it to the clipboard

```
openssl rand -base64 32 | cut -c1-32 | tr -d '\n' | base64 | tr -d '\n' | pbcopy
```
