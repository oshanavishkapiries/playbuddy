@echo off
echo Building PlayBuddy Torrent Search...
go mod tidy
go build -o bin/playbuddy.exe cmd/main.go
if %ERRORLEVEL% EQU 0 (
    echo Build successful!
    echo Run with: bin\playbuddy.exe
) else (
    echo Build failed!
)
pause 