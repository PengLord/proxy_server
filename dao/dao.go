package dao

import (
	"proxy_server/pojo"
	"proxy_server/config"
)



//保存操作
func Save(proxy pojo.ProxyPool)  {
	db := config.GetDb()
    db.Save(&proxy)
}
//查询所有
func FindAll() []pojo.ProxyPool{
	db:=config.GetDb()
	var po []pojo.ProxyPool
	db.Find(&po)
	return po
}

func DeleteById(id string)  {
	db := config.GetDb()
	db.Model(&pojo.ProxyPool{}).Delete(pojo.ProxyPool{Id:id})
}

func SvaeDomain(daomain pojo.Domain)  {
	db := config.GetDb()
	db.Create(&daomain)
}

func FindByDomain(domain string) pojo.Domain {
	db:=config.GetDb()
	po := pojo.Domain{}
	db.Where("domain=?",domain).First(&po)
	return po
}
