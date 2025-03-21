#!/usr/bin/env bash

set -e

if [ -z "$(which mkcert)" ]; then
  echo "mkcert not installed. Installing....."
  HOMEBREW_NO_AUTO_UPDATE=1 brew install mkcert
fi

check_and_generate_certificate() {
  if [ ! -f "$(mkcert -CAROOT)/rootCA.pem" ]; then
    cat << EOF
#####################################################################
# ERROR: The mkcert root CA certificate has not been generated.     #
# Please run the following command to install it:                   #
#                                                                   #
#   mkcert --install                                                #
#                                                                   #
# This is required for nats to start properly.                      #
#####################################################################
EOF
    return 1
  fi


  # Check if server.pem exists and if it is a directory
  if [ -d ./server.pem ]; then
    echo "server.pem is a directory. Deleting certs contents."
    rm -rf ./*
    echo "Directory cleaned up."
  fi

  # Check if the server.pem file exists
  if [ ! -f ./server.pem ]; then
    echo "server.pem does not exist. Generating a new certificate."
    mkcert -cert-file server.pem -key-file server-key.pem localhost 127.0.0.1 ::1
    return 0
  fi

  # Get the certificate expiration date
  enddate=$(openssl x509 -enddate -noout -in ./server.pem | cut -d= -f2)

  # Convert the expiration date to a Unix timestamp
  enddate_timestamp=$(date -j -f "%b %d %T %Y %Z" "$enddate" "+%Y%m%d%H%M%S")

  # Get today's date as a Unix timestamp
  current_timestamp=$(date "+%Y%m%d%H%M%S")

  # Compare the dates
  if [ "$enddate_timestamp" -le "$current_timestamp" ]; then
    echo "Certificate has expired or is expiring today. Generating a new one."
    mkcert -cert-file server.pem -key-file server-key.pem localhost 127.0.0.1 ::1
    return 0
  fi

  echo "Certificate is still valid."
}

mkdir -p ./bin/serve/docker/nats/certs

pushd ./bin/serve/docker/nats/certs

  check_and_generate_certificate

popd

