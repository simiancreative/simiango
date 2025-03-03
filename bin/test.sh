#!/bin/bash

export GOEXPERIMENT=nocoverageredesign

go install github.com/jstemmer/go-junit-report/v2@latest

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors | grep -v examples`

go test -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out 2>&1 $list > test.log

cat test.log

cat test.log | go-junit-report -set-exit-code > junit.xml

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

echo
echo coverage: $coverage
