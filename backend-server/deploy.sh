#!/bin/zsh
GOOS=linux GOARCH=amd64 go build -o ropc
sudo chmod +x ropc

# stop service
ssh root@139.162.170.159 'systemctl stop ropc'

scp ropc root@139.162.170.159:/root/ropc
ssh root@139.162.170.159 'systemctl start ropc'
