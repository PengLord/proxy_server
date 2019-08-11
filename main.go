package main

import (
	"proxy_server/utlis"
	"proxy_server/dao"
	"proxy_server/config"
	"proxy_server/pojo"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/rand"
	"proxy_server/amqphelper"
	"sync"
)

func main() {
	//加载配置文件
	con := config.New("./config/config.yml")
	config.NewDb(con)
	//初始化消息队列
	mq := amqphelper.New(&amqphelper.Config{
		AmqpUrl: con.AmqUrl,
	})
	var proxyw []pojo.ProxyPool
	var proxyz []pojo.ProxyPool
	var Lock sync.RWMutex

	//查询数据
	go utlis.Runtime(3, func() {
		proxys := dao.FindAll()
		for _, v := range proxys {
			Lock.Lock()
			if v.IsCn == 0 {
				proxyz = append(proxyz, v)
			} else {
				proxyw = append(proxyw, v)
			}
			Lock.Unlock()
		}

	})
	//提供服务
	en := gin.Default()
	en.GET("/proxy_pool/api/proxys", func(context *gin.Context) {
		domain := context.Query("domain")
		//查询域名是否以标记
		po := dao.FindByDomain(domain)
		if po.Id <= 0 {
			//发送消息
			po.Domain = domain
			mq.PushWithObject("", con.WorkerName, 5, po)
		}

		if po.IsCn == false { //国外
			//生成随机数
			Lock.RLock()
			su := rand.Intn(len(proxyw))
			proxy := proxyw[su].Proxy
			Lock.RUnlock()
			context.JSON(200, &Resp{
				Msg: "success",
				Data: Data{
					Proxy: proxy,
				},
			})
		} else {

			//国内
			//生成随机数
			Lock.RLock()
			su := rand.Intn(len(proxyz))
			proxy := proxyw[su].Proxy
			Lock.RUnlock()
			context.JSON(200, &Resp{
				Msg: "success",
				Data: Data{
					Proxy: proxy,
				},
			})
		}
	})
	en.Run(":8080")
}

type Resp struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

type Data struct {
	Proxy string `json:"proxy"`
}
