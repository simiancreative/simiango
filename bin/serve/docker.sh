#!/usr/bin/env bash

./bin/serve/gen-certs.sh

pushd ./bin/serve/docker

make start

popd

