#!/usr/bin/env bash

pushd ./examples
go mod tidy
popd

air -c "./bin/serve/air.toml"
