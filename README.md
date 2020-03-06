# 这个是一个博客程序

## 从源代码安装

### 代码下载
```bash
git clone https://github.com/buexplain/go-blog.git
```

### 准备配置文件
```bash
cd go-blog && copy config.example.toml config.toml
```

### 初始化数据库
```bash
go build -o artisan.exe artisan.go && artisan.exe db import -f database/init.sql
```

### 跑起来
```bash
go run main.go
```

### 进入博客后台
1. 后台地址 http://localhost:1991/backend/sign
2. 账号 admin
3. 密码 123456

## 二次开发相关的命令

### 同步models结构体到数据库
```bash
go build -o artisan.exe artisan.go && artisan.exe db sync
```

### 制作 `database/init.sql`
```bash
go build -o artisan.exe artisan.go && artisan.exe db dump -f database/init.sql
```

### 静态文件打包
```bash
go build -o artisan.exe artisan.go && artisan.exe asset pack
```

### 引用本地包
```bash
go mod edit -replace=github.com/buexplain/go-fool=F:/go-fool
go mod edit -replace=github.com/buexplain/go-fool=C:\Edisk\code\go-fool
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
```
