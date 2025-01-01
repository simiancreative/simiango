#!/bin/bash

set -e

export GOEXPERIMENT=nocoverageredesign

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors | grep -v examples`

go test -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out $list
go tool cover -func coverage.out > coverage-report.txt

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

echo
echo coverage: $coverage
