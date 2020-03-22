#!/bin/bash
# 如果提示 /bin/bash^M: 坏的解释器: 没有那个文件或目录，则是因为结束符的不通，windows下是\r\n，linux下是\n的原因
# sed -i 's/\r$//' installer-linux.sh

function isError() {
  if [ $? != 0 ]; then
      exit $?
  fi
}

# 判断是否已经安装过了
if [ -f "installed.lock" ]; then
  echo "already installed, please execute: blog.bin"
  exit 0
fi

echo "install start..."

# 生成配置文件
cp config.example.toml config.toml
isError

# 释放静态文件
./artisan.bin asset unpack
isError

# 同步表结构到数据库
./artisan.bin db sync
isError

# 导入database/init.sql文件
./artisan.bin db import -f ./database/init.sql
isError

# 启动程序
if [ $? == 0 ]; then
  echo "install successfully, please execute: blog.bin"
  date -d today +"%Y-%m-%d %H/%M/%S" > installed.lock
else
    echo "install failed"
fi