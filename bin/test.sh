#!/bin/bash

set -e

go install github.com/mattn/goveralls@latest

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors`

go test -cover -coverprofile=coverage.out $list
goveralls -coverprofile=coverage.out -repotoken=${COVERALLS_TOKEN}

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

echo
echo coverage: $coverage
