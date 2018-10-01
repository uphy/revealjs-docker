#!/bin/sh

if [ $# == 0 ]; then
  echo "Specify the version. e.g., 3.7.0"
  exit 1
fi

sed -i "" -e 's@uphy\/reveal.js:[0-9\.]*@uphy\/reveal.js:'"$1@" README.md
sed -i "" -e 's@VERSION=.*@VERSION='"$1"'@' Dockerfile
git add .
git commit -m "Updated reveal.js version."
git tag -am "Updated reveal.js version." $1