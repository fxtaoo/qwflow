package main

import (
	"log"
	"qwflow/conf"
	"qwflow/web"
)

func main() {
	var conf conf.Conf
	// 初始化数据
	err := conf.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer conf.Mysql.DB.Close()

	web.Start(&conf)
}
