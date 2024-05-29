#!/bin/bash

version=$1

# Version must be in the format x.y.z
if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Version must be in the format x.y.z"
  exit 1
fi

cd $(dirname $0)/..
sed -i '' -e "s/uphy\/reveal\.js:[0-9]*\.[0-9]*\.[0-9]*/uphy\/reveal.js:$version/g" README.md
sed -i '' -e "s/ARG REVEALJS_VERSION=.*/ARG REVEALJS_VERSION=$version/" Dockerfile