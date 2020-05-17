package config

//业务相关的配置
type Business struct {
	//超级角色，该角色会自动拥有所有的权限节点
	SuperRoleID int
	//附件上传配置
	Upload upload
	//三方登录
	OAuth map[string]OAuth
}

type OAuth struct {
	ID string
	Secret string
	RedirectUri string
}

type upload struct {
	//允许的后缀与媒体类型
	Accept map[string]string
	//保存的路径
	Save string
}

func (this upload) Ext() []string {
	str := make([]string, 0, len(this.Accept))
	for k,_ := range this.Accept  {
		str = append(str, string(k))
	}
	return str
}

func (this upload) MimeType() []string {
	m := map[string]bool{}
	i := 0
	for _, v := range this.Accept  {
		tmp := string(v)
		if tmp == "" {
			continue
		}
		m[string(v)] = true
		i++
	}
	str := make([]string, 0, i)
	for k, _ := range m {
		str = append(str, k)
	}
	return str
}
