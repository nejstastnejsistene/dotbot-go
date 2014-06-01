#!/bin/sh
CC=arm-linux-androideabi-gcc \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
CGO_ENABLED=1 \
GOROOT=/home/peter/code/go \
go build -tags android
