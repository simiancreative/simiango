#!/bin/bash

export GOEXPERIMENT=nocoverageredesign

go install github.com/jstemmer/go-junit-report/v2@latest

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors | grep -v examples`

go test -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out $list 2>&1 | tee test.log
test_exit=${PIPESTATUS[0]}

cat test.log | go-junit-report > junit.xml

if [ $test_exit -ne 0 ]; then
  exit $test_exit
fi

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

echo
echo coverage: $coverage
