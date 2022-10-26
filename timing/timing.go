package timing

import (
	"log"
	"qwflow/conf"
	"qwflow/wangsu"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func Start() {
	// 定时运行
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc("0 3 * * *", func() { YesterdayFlow() })
	c.Start()

	wg.Wait()

	// YesterdayFlow()

}

func YesterdayFlow() {
	var conf conf.Conf
	// 初始化数据
	err := conf.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer conf.Mysql.DB.Close()

	end := time.Now()
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
	err = d.LiveDataChannelPeak(&conf.Wangsu, &conf.Mysql)
	if err != nil {
		log.Fatal(err)
	}

	// 网宿 cdn
	d2 := &wangsu.DateChannelPeak{
		Begin: bengin,
		End:   end,
	}
	err = d2.DataChannelPeak(&conf.Wangsu, "dl-https;download;live-https;web;web-https")
	if err != nil {
		log.Fatal(err)
	}
	d2.Save(&conf.Mysql, "WangsuCdnFlow")

}