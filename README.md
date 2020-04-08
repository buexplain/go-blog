# 这个是一个博客程序

## 从源代码安装

```bash
git clone https://github.com/buexplain/go-blog.git
cd go-blog
chmod u+x ./build-linux.sh
./build-linux.sh
cd ./build
chmod u+x ./installer.sh
./installer.sh
chmod u+x ./blog.bin
./blog.bin
```

## 发布程序

### 编译

```bash
# 编译程序 linux下是 build-linux.sh
# 如果提示 /bin/bash^M: 坏的解释器: 没有那个文件或目录，请执行 sed -i 's/\r$//' build-linux.sh
build-windows.bat
```

### 注意事项
如果你编译失败，例如提示一下错误之一：
* std.ERROR: 导出数据库到文件失败: column LastTime type []uint8 convert to time.Time error
* 编译成功，但是测试`installer.sh`脚本的时候，提示导入数据库失败
* 编译成功，测试`installer.sh`脚本也成功，但是数据库里面的数据少了一部分，比如新添加的文章丢失了
以上错误的原因有可能是`database/database.shm`、`database/database.wal`文件，不能跨平台的原因。
比如你在windows下开发，生成了`database/database.db`、`database/database.shm`、`database/database.wal`文件。
然后你把它们拷贝到linux下，然后你运行了`build-linux.bin`，编译失败了。
此时可以重做一个`database/database/db`数据库文件，请按如下步骤操作：
1. 首先，windows下运行导出sql命令`go build -o artisan.exe artisan.go && artisan.exe db dump -m 64 -f database/init.sql`
2. 然后，将`config.toml`里面`DSN`的值改为`database.db`，去掉所有的连接参数
3. 然后，windows下运行导入sql命令`go build -o artisan.exe artisan.go && artisan.exe db sync && artisan.exe db import -f ./database/init.sql`
4. 最后，将`database/database.db`，拷贝到Linux，执行编译命令`./build-linux.bin`
当然你也可以导出`database/init.sql`后，将其拷贝到Linux，然后再在Linux下生成`database/database.db`，然后再编译。

## 其它

### 二次开发相关命令
```bash
# 导出 database/init.sql
go build -o artisan.exe artisan.go && artisan.exe db dump -m 64 -f database/init.sql
# 导入  database/init.sql
go build -o artisan.exe artisan.go && artisan.exe db sync && artisan.exe db import -f ./database/init.sql

# 静态文件打包
go build -o artisan.exe artisan.go && artisan.exe asset pack
```

### 引用本地包
```bash
go mod edit -replace=github.com/buexplain/go-fool=F:/go-fool
go mod edit -replace=github.com/buexplain/go-fool=C:\Edisk\code\go-fool
go mod edit -replace=github.com/buexplain/go-fool=/mnt/winEdisk/code/go-fool
```

### 包升级到最新版本
```bash
go get -u 包路径@[版本号,保持最新请使用latest 或者 master]
go get -u xorm.io/xorm@latest
go get -u github.com/88250/lute@latest
go get -u github.com/mojocn/base64Captcha@latest
go get -u github.com/spf13/cobra@latest
go get -u github.com/gorilla/sessions@latest
go get -u github.com/kevinburke/go-bindata@latest
go get -u github.com/mattn/go-sqlite3@latest
go get -u github.com/BurntSushi/toml@latest
go get -u github.com/gorilla/securecookie@latest
go get -u github.com/buexplain/go-flog@latest
go get -u github.com/buexplain/go-fool@latest
go get -u github.com/cloudflare/tableflip@latest
```

## License
[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
