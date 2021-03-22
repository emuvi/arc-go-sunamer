#/bin/bash
rm -rf dist
mkdir dist
mkdir dist/linux-amd64
mkdir dist/darwin-amd64
mkdir dist/windows-amd64
GOARCH=amd64
GOOS=linux
go build -o dist/linux-amd64/sunamer
GOOS=darwin
go build -o dist/darwin-amd64/sunamer
GOARCH=386
GOOS=windows
go build -o dist/windows-amd64/sunamer.exe