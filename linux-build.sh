#!/bin/bash
# 如果提示 /bin/bash^M: 坏的解释器: 没有那个文件或目录，则是因为结束符的不通，windows下是\r\n，linux下是\n的原因
# sed -i 's/\r$//' build-linux.sh

function isError() {
  if [ $? != 0 ]; then
      exit $?
  fi
}

clear

echo "build start..."

# 检查编译输出目录
if [ ! -d "build" ]; then
  mkdir build
else
  rm -fr build/*
fi

# 输出readme文件
echo "一、文件介绍" > ./build/readme.txt
echo "  1、artisan.bin 各种命令集合，比如数据导入导出命令等 ./artisan.bin -h 查看详细。" >> ./build/readme.txt
echo "  2、blog.bin 博客主程序。" >> ./build/readme.txt
echo "  3、config.example.toml 配置文件范例。" >> ./build/readme.txt
echo "  4、installer.sh 安装脚本，用来初始化整个应用。" >> ./build/readme.txt
echo "  5、readme.txt 自述文件" >> ./build/readme.txt
echo "二、安装步骤" >> ./build/readme.txt
echo "  1、执行 ./installer.sh 等待一段时间会提示安装成功。" >> ./build/readme.txt
echo "  2、执行 ./blog.bin 不要关闭它，屏幕会提示网址。" >> ./build/readme.txt
echo "  3、打开浏览器，输入 blog.bin 提示的网址，有些云服务器和虚拟机，提示的网址是内网ip，请替换成公网ip。" >> ./build/readme.txt
echo "  4、账号 admin 密码 123456" >> ./build/readme.txt
echo "三、生产环境部署" >> ./build/readme.txt
echo "  nohup ./blog.bin 1>info.log 2>error.log &" >> ./build/readme.txt
echo "四、更多信息" >> ./build/readme.txt
echo "  https://github.com/buexplain/go-blog" >> ./build/readme.txt

#复制配置文件
cp config.example.toml build/config.example.toml
# 复制安装器
cp linux-installer.sh build/installer.sh

# 打开cgo
export CGO_ENABLED=1

# 生成配置文件
if [ ! -f "config.toml" ]; then
  cp config.example.toml config.toml
fi

# 更新依赖包
go mod tidy
isError

# 编译命令行程序
go build -ldflags "-s -w" -o artisan.bin artisan.go
isError

# 打包静态文件
./artisan.bin asset pack
isError

# 删除artisan.bin
rm -fr ./artisan.bin
isError

# 编译发布程序
go build -ldflags "-s -w" -o ./build/blog.bin main.go && go build -ldflags "-s -w" -o ./build/artisan.bin artisan.go
if [ $? == 0 ]; then
  echo "build successfully"
  cd build && pwd && ls
else
    echo "build failed"
fi