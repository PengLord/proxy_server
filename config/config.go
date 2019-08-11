package config

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/jinzhu/gorm"
)

type Config struct {
	Db string		`yaml:"db"`
	ProxyUrl string	 `yaml:"proxyurl"`
	ValidTime int	  `yaml:"validtime"`
	AmqUrl string 		`yaml:"amqurl"`
	WorkerName	 string		`yaml:"queuename"`
}

func New(path string) Config {
	con := Config{}
	buff,err := ioutil.ReadFile(path)
	if err!= nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(buff,&con)
	return con
}

var db *gorm.DB
func NewDb(config Config)  {
	var err error
	db,err = gorm.Open("mysql",config.Db)
	if err!=nil{
		log.Fatal(err)
	}
}

func GetDb() *gorm.DB {
	return db
}
