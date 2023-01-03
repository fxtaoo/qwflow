package timing

import (
	"log"
	"qwflow/conf"
	"qwflow/wangsu"
	"time"

	"github.com/robfig/cron/v3"
)

func Start() {
	var conf conf.Conf
	// 初始化数据
	err := conf.Init()
	if err != nil {
		log.Fatal(err)
	}

	// 定时运行
	c := cron.New()
	// 获取昨天数据
	c.AddFunc("0 3 * * *", func() {
		// 数据库初始化
		conf.Mysql.Init()
		defer conf.Mysql.DB.Close()

		now := time.Now()

		// 获取昨天七牛网宿相关数据
		err = YesterdayFlow(&conf, now)
		if err != nil {
			log.Fatal(err)
		}

		// 流量日环比增幅超过设定值邮件告警
		err = conf.Alerts.Calc(&conf.Mysql)
		if err != nil {
			log.Fatal(err)
		}
		conf.Alerts.SendMail()
	})
	// 周一发送图片流量报表
	// 图片需要提前生成好
	c.AddFunc("0 5 * * 1", func() {
		if conf.ChartMail.Switch {
			conf.ChartMail.SendMail("live", "d14")
			conf.ChartMail.SendMail("cdn", "d14")
		}
	})

	c.Start()
	select {}
}

func YesterdayFlow(conf *conf.Conf, end time.Time) error {
	bengin := end.AddDate(0, 0, -1)

	// 七牛直播
	conf.Qiniu.Pili.HubsFlow.DataFlows(conf.Qiniu.Pili.Manager, &conf.Mysql, bengin, end)

	// 七牛 cdn
	conf.Qiniu.Cdn.CndFlows.DataFlows(conf.Qiniu.Cdn.Manager, &conf.Mysql, bengin, end)

	// 网宿直播
	d := &wangsu.DateChannelPeak{
		Begin: bengin,
		End:   end,
	}
	err := d.LiveDataChannelPeak(&conf.Wangsu, &conf.Mysql)
	if err != nil {
		return err
	}

	// 网宿 cdn
	d2 := &wangsu.DateChannelPeak{
		Begin: bengin,
		End:   end,
	}
	err = d2.DataChannelPeak(&conf.Wangsu, "dl-https;download;live-https;web;web-https")
	if err != nil {
		return err
	}
	d2.Save(&conf.Mysql, "WangsuCdnFlow")

	return nil
}
