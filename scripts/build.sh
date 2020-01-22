#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-s -w' -o server cmd/server.go && sudo docker build -f build/package/Dockerfile -t bingo:latest .
#rm server


