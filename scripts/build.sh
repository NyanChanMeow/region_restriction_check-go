#!/bin/bash

rm -rf build/*

export CGO_ENABLED="0"
export GOROOT_FINAL="/usr"

echo Building Linux amd64
export GOOS=linux
export GOARCH=amd64
go build -a -trimpath -asmflags "-w -s" -ldflags "-w -s" -o build/rrc-linux-amd64 || exit $?

upx build/*

cp config.example.json build/config.example.json