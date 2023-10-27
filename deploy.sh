#!/bin/zsh
GOOS=linux GOARCH=amd64 go build -o ropc
docker build -t unitz007/ropc .
docker push unitz007/ropc
ssh root@143.42.56.143 'kubectl delete -f ropc/ropc-deployment.yaml && kubectl apply -f ropc/ropc-deployment.yaml'
