#!/usr/bin/env bash

pushd ./examples
export PORT=8081
go run . -env dev
popd

