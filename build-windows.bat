@echo off

cls

echo build start...

rem 检查编译输出目录
if not exist build (
    mkdir build
) else (
    if exist build\database (
        rem 删除文件夹下的文件以及文件夹本身
        rd /s /q build\database
    )
    if exist build\resources (
        rem 删除文件夹下的文件以及文件夹本身
        rd /s /q build\resources
    )
    if exist build\storage (
        rem 删除文件夹下的文件以及文件夹本身
        rd /s /q build\storage
    )
    if exist build\uploads (
        rem 删除文件夹下的文件以及文件夹本身
        rd /s /q build\uploads
    )
    rem 删除文件夹下的文件
    del /f /s /q build\*
)

rem 复制配置文件
copy config.example.toml .\build\config.example.toml /A/Y

rem 复制安装器
copy installer-windows.bat .\build\installer-windows.bat /A/Y

rem 打开cgo
SET CGO_ENABLED=1

rem 更新依赖包
go mod tidy
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 编译命令行程序
go build -ldflags "-s -w" -o artisan.exe artisan.go
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 同步表结构到数据库
artisan.exe db sync
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 导出表数据
artisan.exe db dump -m 64 -f database/init.sql
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 打包静态文件
artisan.exe asset pack
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 复制artisan.exe
move /Y artisan.exe .\build\artisan.exe
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 编译发布程序
go build -ldflags "-s -w" -o ./build/blog.exe main.go
if %errorlevel% == 0 (
    echo build successfully
    dir build
) else (
    echo build failed
)
pause