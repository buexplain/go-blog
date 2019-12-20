@echo off
echo build start...
if not exist build (
    mkdir build
)
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w" -o ./build/blog-amd64.exe main.go && go build -ldflags "-s -w" -o ./build/artisan-amd64.exe artisan.go
if %errorlevel% == 0 (
echo build successfully
) else (
echo build failed
)
pause