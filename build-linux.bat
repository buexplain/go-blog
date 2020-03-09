@echo off
echo build start...
if not exist build (
    mkdir build
)
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go mod tidy
go build -ldflags "-s -w" -o ./build/blog-amd64.bin main.go && go build -ldflags "-s -w" -o ./build/artisan-amd64.bin artisan.go
if %errorlevel% == 0 (
echo build successfully
) else (
echo build failed
)
pause