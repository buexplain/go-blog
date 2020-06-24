@echo off

rem 批处理教程 http://bbs.bathome.net/thread-18-1-1.html
rem 批处理常用命令及用法大全 http://bbs.bathome.net/thread-39-1-1.html

rem 声明采用UTF-8编码
chcp 65001

rem 清空屏幕
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

rem 输出readme文件
echo 一、文件介绍 > build\readme.txt
echo   1、artisan.exe 各种命令集合，比如数据导入导出命令等 artisan.exe -h 查看详细。 >> build\readme.txt
echo   2、blog.exe 博客主程序。 >> build\readme.txt
echo   3、config.example.toml 配置文件范例。 >> build\readme.txt
echo   4、installer.bat 安装脚本，用来初始化整个应用。 >> build\readme.txt
echo   5、readme.txt 自述文件 >> build\readme.txt
echo 二、安装步骤 >> build\readme.txt
echo   1、双击 installer.bat 等待一段时间会提示安装成功。 >> build\readme.txt
echo   2、双击 blog.exe 不要关闭它，屏幕会提示网址。 >> build\readme.txt
echo   3、打开浏览器，输入 blog.exe 提示的网址，有些云服务器和虚拟机，提示的网址是内网ip，请替换成公网ip。 >> build\readme.txt
echo   4、账号 admin 密码 123456 >> build\readme.txt
echo 三、更多信息 >> build\readme.txt
echo   https://github.com/buexplain/go-blog >> build\readme.txt

rem 复制配置文件
copy config.example.toml .\build\config.example.toml /A/Y

rem 复制安装器
copy installer-windows.bat .\build\installer.bat /A/Y

rem 打开cgo
SET CGO_ENABLED=1

rem 生成配置文件
if not exist config.toml (
    copy config.example.toml config.toml
)

rem 更新依赖包
go mod tidy
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 编译命令行程序
go build -ldflags "-s -w" -o artisan.exe artisan.go
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 打包静态文件
artisan.exe asset pack
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 删除artisan.exe
del /f /s /q artisan.exe
if %errorlevel% NEQ 0 exit /b %errorlevel%

rem 编译发布程序
go build -ldflags "-s -w" -o ./build/blog.exe main.go && go build -ldflags "-s -w" -o ./build/artisan.exe artisan.go
if %errorlevel% == 0 (
    echo build successfully
    dir build
) else (
    echo build failed
)
pause