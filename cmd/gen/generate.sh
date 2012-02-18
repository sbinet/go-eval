#!/bin/bash
go build . || exit 1
./gen > expr1.go
gofmt -w expr1.go
