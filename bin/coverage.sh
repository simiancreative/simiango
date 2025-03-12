#!/usr/bin/env bash

set -e

# Default values
should_open=false
package_path="./..."

# Function to display usage
usage() {
    echo "Usage: $0 [-p|--package <package_path>] [-o|--open] [-h|--help]"
    echo "  -p, --package   Specify the package path to test (default: all packages)"
    echo "  -o, --open      Open the coverage report in the default browser"
    echo "  -h, --help      Display this help message"
}

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -p|--package) package_path="$2"; shift ;;
        -o|--open) should_open=true ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown parameter passed: $1"
            usage
            exit 1
            ;;
    esac
    shift
done

# TODO: work through lint issues and enable golanci-lint
# lint
# golangci-lint run $package_path

# unit test
go test $package_path -coverpkg=$package_path -race -covermode=atomic -coverprofile=coverage.out
go tool cover -html=coverage.out -o ~/Desktop/coverage.html

# report
echo ""
echo "==="
echo "Coverage report is generated at ~/Desktop/coverage.html"
echo "==="

# open
if [ "$should_open" = true ]; then
    open ~/Desktop/coverage.html
fi
