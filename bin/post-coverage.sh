#!/bin/bash

set -e

total=$1

if [ -z "$GITHUB_SHA" ]
then
  echo no sha
  exit 0
fi

curl \
  --header "Authorization: Token e795b6b9-ce88-485c-bd99-6a81ae66fc40" \
  --header "Content-Type: application/json" \
  --data "{
    \"value\":\"${total}\",
    \"sha\":\"${GITHUB_SHA}\"
  }" \
  https://seriesci.com/api/simiancreative/simiango/coverage/one
