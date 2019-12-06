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
go build -o artisan.exe artisan.go && artisan.exe db sync && artisan.exe db import -f database/init.sql
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

### 初始化一个管理员用户
```bash
go build -o artisan.exe artisan.go && artisan.exe db sync &&  artisan.exe db addAdmin -a admin -p 123456
```

### `database/init.sql`制作
下面的命令是导出数据，表结构与表索引不做导出。
```bash
go build -o artisan.exe artisan.go && artisan.exe db dump -m 64 -f database/init.sql
```

### 静态文件打包
```bash
go build -o artisan.exe artisan.go && artisan.exe asset pack
```

### 引用本地包
```bash
go mod edit -replace=github.com/buexplain/go-fool=F:/go-fool
```
