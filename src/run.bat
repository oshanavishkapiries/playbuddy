@echo off
echo Starting PlayBuddy Torrent Search...
go mod tidy
go run cmd/main.go
pause 