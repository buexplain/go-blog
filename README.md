# 这个是一个博客程序

## 从源代码安装

```bash
# 代码下载
git clone https://github.com/buexplain/go-blog.git
# 准备配置文件
cd go-blog && copy config.example.toml config.toml
# 初始化数据库，初始的账号密码是 admin，123456
go build -o artisan.exe artisan.go && artisan.exe db import -f database/init.sql
# 编译运行
go run main.go
```

## 发布程序

```bash
# 编译程序 linux下是 build-linux.bin
build-windows.bat
```

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
```
