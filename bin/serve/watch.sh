#!/usr/bin/env bash

fswatch ./ -e ".*" -i "\\.go$" | xargs -n1 -I{} overmind restart examples
