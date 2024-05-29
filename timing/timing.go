package timing

import (
	"fmt"
	"log"
	"qwflow/conf"
	"qwflow/wangsu"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func Start() {
	var Conf conf.Conf
	// 初始化数据
	err := Conf.Init()
	if err != nil {
		log.Fatal(err)
	}

	// 定时运行
	c := cron.New()
	// 获取昨天数据
	c.AddFunc(Conf.TimingCron.GetDayFlow, func() {
		var conf conf.Conf
		// 初始化数据
		err := conf.Init()
		if err != nil {
			log.Fatal(err)
		}
		// 数据库初始化
		conf.Mysql.Init()
		defer conf.Mysql.DB.Close()

		end := time.Now()
		bengin := end.AddDate(0, 0, -1)

		// 获取昨天七牛网宿相关数据
		err = GetDayFlow(&conf, bengin, end, "qiniuwangsu", "livecdn", false)
		if err != nil {
			log.Fatal(err)
		}
	})

	// 流量日环比增幅超过设定值邮件告警
	c.AddFunc(Conf.TimingCron.Alerts, func() {
		var conf conf.Conf
		// 初始化数据
		err := conf.Init()
		if err != nil {
			log.Fatal(err)
		}
		// 数据库初始化
		conf.Mysql.Init()
		defer conf.Mysql.DB.Close()

		// 流量日环比增幅超过设定值邮件告警
		err = conf.Alerts.Calc(&conf.Mysql)
		if err != nil {
			log.Fatal(err)
		}
		conf.Alerts.SendMail()
	})

	// 周一发送图片流量报表
	// 图片需要提前生成好
	c.AddFunc(Conf.TimingCron.ChartMail, func() {
		var conf conf.Conf
		// 初始化数据
		err := conf.Init()
		if err != nil {
			log.Fatal(err)
		}

		if conf.ChartMail.Switch {
			conf.ChartMail.SendMail("live", conf.ChartMail.ImgName)
			conf.ChartMail.SendMail("cdn", conf.ChartMail.ImgName)
		}
	})

	c.Start()
	select {}
}

func GetDayFlow(conf *conf.Conf, bengin, end time.Time, qiniuwangsu, livecdn string, delete bool) error {
	// 删除数据
	dataDel := func(table string) {
		if delete {
			sql := fmt.Sprintf("DELETE FROM %s WHERE date=\"%s\"", table, bengin.Format("2006-01-02"))
			stmt, _ := conf.Mysql.GetStmt(sql)
			stmt.Exec()
			defer stmt.Close()
		}
	}
	if strings.Contains(qiniuwangsu, "qiniu") {
		if strings.Contains(livecdn, "live") {
			// 七牛直播
			dataDel("QiniuHubsFlow")
			conf.Qiniu.Pili.HubsFlow.DataFlows(conf.Qiniu.Pili.Manager, &conf.Mysql, bengin, end)
		}
		if strings.Contains(livecdn, "cdn") {
			// 七牛 cdn
			dataDel("QiniuCdnsFlow")
			conf.Qiniu.Cdn.CndFlows.DataFlows(conf.Qiniu.Cdn.Manager, &conf.Mysql, bengin, end)
		}
	}

	if strings.Contains(qiniuwangsu, "wangsu") {
		if strings.Contains(livecdn, "live") {
			// 网宿直播
			dataDel("WangsuLiveFlow")
			d := &wangsu.DateChannelPeak{
				Begin: bengin,
				End:   end,
			}
			err := d.LiveDataChannelPeak(&conf.Wangsu, &conf.Mysql)
			if err != nil {
				return err
			}

		}

		if strings.Contains(livecdn, "cdn") {
			/// 网宿 cdn
			dataDel("WangsuCdnFlow")
			d2 := &wangsu.DateChannelPeak{
				Begin: bengin,
				End:   end,
			}
			err := d2.DataChannelPeak(&conf.Wangsu, "dl-https;download;live-https;web;web-https")
			if err != nil {
				return err
			}
			d2.Save(&conf.Mysql, "WangsuCdnFlow")
		}

	}

	return nil
}
