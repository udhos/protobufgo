#!/bin/bash

go install golang.org/x/vuln/cmd/govulncheck@latest
go install golang.org/x/tools/cmd/deadcode@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

protoc -I=. --go_out=proto proto/addressbook.proto

gofmt -s -w .

revive ./...

gocyclo -over 15 .

go mod tidy

govulncheck ./...

deadcode ./cmd/*

go env -w CGO_ENABLED=1

go test -race ./...

#go test -bench=BenchmarkController ./cmd/kubecache

go env -w CGO_ENABLED=0

go install ./...

go env -u CGO_ENABLED
