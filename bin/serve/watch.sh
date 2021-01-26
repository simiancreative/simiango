#!/usr/bin/env bash

fswatch ./modules/auth -e ".*" -i "\\.go$" | xargs -n1 -I{} overmind restart auth
