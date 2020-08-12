#!/bin/sh -eu

rm -rf dist
mkdir dist
go build -o dist/main main.go
