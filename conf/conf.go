package conf

import (
	"encoding/json"
	"os"
	"qwflow/alert"
	"qwflow/chartmail"
	"qwflow/mysql"
	"qwflow/qiniu"
	"qwflow/wangsu"

	"github.com/gin-gonic/gin"
)

type Conf struct {
	Mysql       mysql.Mysql         `json:"mysql"`
	Qiniu       qiniu.Qiniu         `json:"qiniu"`
	Wangsu      wangsu.WangSu       `json:"wangsu"`
	Accounts    gin.Accounts        `json:"accounts"`
	Alerts      alert.Alerts        `json:"alerts"`
	ChartMail   chartmail.ChartMail `json:"chartMail"`
	CdnOtherGB  int                 `json:"cdnOtherGB"`
	LiveOtherTB float64             `json:"liveOtherTB"`
	TimingCron  struct {
		GetDayFlow string `json:"GetDayFlow"`
		Alerts     string `json:"Alerts"`
		ChartMail  string `json:"ChartMail"`
	} `json:"TimingCron"`
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
	// 七牛初始化
	c.Qiniu.Init()
	// 网宿初始化
	c.Wangsu.Init()
	// 流量日环比增幅超过设定值邮件告警
	c.Alerts.Init()

	return nil
}
