#!/bin/sh
CC=arm-linux-androideabi-gcc \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
CGO_ENABLED=1 \
go build -tags android
