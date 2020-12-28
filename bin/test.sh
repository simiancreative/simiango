#!/bin/bash

set -e

list=`go list ./... | grep -v mocks | grep -v docs | grep -v errors`

go test -cover $list

coverage=$(go tool cover -func=coverage.out)

echo
echo coverage: $(echo $coverage | grep total | awk '{print $3}')
