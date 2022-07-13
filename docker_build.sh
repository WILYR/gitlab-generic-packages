#!/bin/bash

CGO_ENABLED=0 go build -o main.sh main.go
docker build --network=host --label=gitlab-generic-packages --label=v0.1.5 -t gitlab-generic-packages:0.1.5 .
rm main.sh