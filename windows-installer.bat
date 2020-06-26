@echo off

rem 判断是否已经安装过了
if exist installed.lock (
    echo already installed, please execute: blog.exe
    exit /b 0
)

echo install start...

rem 生成配置文件
if not exist config.toml (
    copy config.example.toml config.toml /A/Y
    if %errorlevel% NEQ 0 exit /b %errorlevel%
)

rem 释放静态文件
artisan.exe asset unpack
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 导入database/init.sql文件
artisan.exe db importInit
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 启动程序
if %errorlevel% == 0 (
    echo install successfully, please execute: blog.exe
    rem 写入安装锁
    echo %DATE:~0,4%-%DATE:~5,2%-%DATE:~8,2% %TIME:~0,2%/%TIME:~3,2%/%TIME:~6,2% > installed.lock
) else (
    echo install failed
)
pause