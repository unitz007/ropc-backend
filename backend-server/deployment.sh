#!/usr/bin/env bash

# sets gin to release mode
export GIN_MODE=release

# run tests
if go test ./...; then
  echo "test passed"
else
  echo "Test failed. Exiting..."
  exit 0
fi

# build go project to executable
GOOS="linux" GOARCH="amd64" go build .

# build container image
docker build -t backend-server .

# deploy image
docker-compose up -d
