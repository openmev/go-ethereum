#!/bin/bash
go mod vendor
go build
make geth-linux-amd64
make geth-darwin-amd64
