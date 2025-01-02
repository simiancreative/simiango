#!/bin/bash

set -e

export GOEXPERIMENT=nocoverageredesign

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors | grep -v examples`

go test -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out 2>&1 $list | go-junit-report -set-exit-code > report.xml

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

echo
echo coverage: $coverage
