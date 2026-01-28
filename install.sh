#! /usr/bin/env bash

go build -o pilot main.go
sudo mv pilot /usr/local/bin

echo "Successfully installed to /usr/local/bin/pilot"
