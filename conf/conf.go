package conf

import (
	"encoding/json"
	"os"
	"qwflow/mysql"
	"qwflow/qiniu"
	"qwflow/wangsu"
)

type Conf struct {
	Mysql  mysql.Mysql   `json:"mysql"`
	Qiniu  qiniu.Qiniu   `json:"qiniu"`
	Wangsu wangsu.WangSu `json:"wangsu"`
}

func (c *Conf) Init() error {
	// 配置从文件读取
	confFile, err := os.ReadFile("./conf.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(confFile, &c)
	if err != nil {
		return err
	}

	// 数据库初始化
	c.Mysql.Init()
	// 七牛初始化
	c.Qiniu.Init()
	// 网宿初始化
	c.Wangsu.Init()

	return nil
}
