package config

//业务相关的配置
type Database struct {
	DSN string
	//设置连接池中的最大闲置连接数
	MaxIdleConns int
	// 设置与数据库建立连接的最大数目
	MaxOpenConns int
	//设置缓存，如果为正数，则表示cache设置为多少页（默认页大小为1KB），如果为负数则表示设置为多少KB
	CacheSize int
	//控制台显示日志
	ShowSQL bool
}
