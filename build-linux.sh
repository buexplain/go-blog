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

# 生成配置文件
if [ ! -f "config.toml" ]; then
  cp config.example.toml config.toml
fi

# 输出readme文件
echo "安装步骤" > ./build/readme.txt
echo "1、执行 ./installer-linux.sh 等待一段时间会提示安装成功。" >> ./build/readme.txt
echo "2、执行 ./blog.bin 不要关闭它，屏幕会提示网址。" >> ./build/readme.txt
echo "3、打开浏览器，输入 blog.bin 提示的网址。" >> ./build/readme.txt
echo "4、账号 admin 密码 123456" >> ./build/readme.txt

#复制配置文件
cp config.example.toml build/config.example.toml
# 复制安装器
cp installer-linux.sh build/installer-linux.sh

# 打开cgo
export CGO_ENABLED=1

# 更新依赖包
go mod tidy
isError

# 编译命令行程序
go build -ldflags "-s -w" -o artisan.bin artisan.go
isError

if [ -f "./database/database.db" ]; then
  # 同步表结构到数据库
  ./artisan.bin db sync
  isError

  # 导出表数据
  ./artisan.bin db dump -m 64 -f database/init.sql
  isError
fi

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
  ls build
else
    echo "build failed"
fi