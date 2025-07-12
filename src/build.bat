@echo off
go mod tidy
go build -o bin/playbuddy.exe cmd/main.go