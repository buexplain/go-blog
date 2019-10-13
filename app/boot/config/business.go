package config

//业务相关的配置
type Business struct {
	//附件上传配置
	Upload struct {
		//允许的后缀
		Ext []string
		//保存的路径
		Save string
	}
}

