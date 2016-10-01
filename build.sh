#!/bin/bash

## clean
rm -rf build/

## for Mac
(
  GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/db cmd/db/*.go
  cd ./build/darwin-amd64
  tar cfz db-darwin-amd64.tar.gz db
)

## for windows
(
  GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/db.exe cmd/db/*.go
  cd ./build/windows-amd64
  zip -q db-windows-amd64.tar.gz db.exe
)

## for linux
(
  GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/db cmd/db/*.go
  cd ./build/linux-amd64
  tar cfz db-linux-amd64.tar.gz db
)
