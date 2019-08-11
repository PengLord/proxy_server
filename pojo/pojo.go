package pojo

type ProxyPool struct {
	Id         string `gorm:"primary_key"` //主键
	Proxy      string                      //代理IP
	IsLongTime int                         //有效时间
	IsCn       int                         //国家
	CreatedAt  int64                       //创建时间
}

type Domain struct {
	Id        int `gorm:"primary_key"` //主键
	Domain    string
	IsCn      bool
	CreatedAt int64
}
