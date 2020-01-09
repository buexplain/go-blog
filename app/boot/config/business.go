package config

//业务相关的配置
type Business struct {
	//超级角色，该角色会自动拥有所有的权限节点
	SuperRoleID int
	//附件上传配置
	Upload struct {
		//允许的后缀
		Ext []string
		//input file 选择文件的 accept
		MimeTypes []string
		//保存的路径
		Save string
	}
}
