#!/bin/bash

set -e

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors`

go test -cover -coverprofile=coverage.out $list

coverage=$(go tool cover -func coverage.out | grep total | awk '{print $3}')

#./bin/post-coverage.sh $coverage

echo
echo coverage: $coverage
