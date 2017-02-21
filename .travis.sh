#!/usr/bin/env bash

set -e

echo "-----------------------"
echo "Build examples"
echo "-----------------------"
find -name main.go | grep -v 'vendor' | while read example; do
	pushd `dirname $example` > /dev/null 2>&1
	echo "== Building $example"
	go get -t -v ./...
	go build
	popd > /dev/null 2>&1
done

echo "-----------------------"
echo "Tests and code-coverage"
echo "-----------------------"
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
