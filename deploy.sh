#!/bin/zsh
GOOS=linux GOARCH=amd64 go build -o ropc
docker build -t unitz007/ropc .
docker push unitz007/ropc
