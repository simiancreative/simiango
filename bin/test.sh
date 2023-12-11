#!/bin/bash

set -e

export GOEXPERIMENT=nocoverageredesign

go install github.com/mattn/goveralls@latest

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors | grep -v examples`

go test -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out $list
goveralls -coverprofile=coverage.out -repotoken=${COVERALLS_TOKEN} \
  || echo 'not posted to coveralls'

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

echo
echo coverage: $coverage
