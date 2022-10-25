package web

import (
	"encoding/json"
	"html/template"
	"qwflow/conf"
	"qwflow/echarts"
	"time"

	"github.com/gin-gonic/gin"
)

type WebValue struct {
	Begin time.Time
	End   time.Time

	QiniuLiveLineStack      *echarts.LineStack
	QiniuLiveLineStackFlows *echarts.LineStackFlows

	QiniuLivePie      *echarts.Pie
	QiniuLivePieFlows *echarts.Pie

	QiniuCdnPie      *echarts.Pie
	QiniuCdnPieFlows *echarts.Pie

	QiniuCdnLineStack      *echarts.LineStack
	QiniuCdnLineStackFlows *echarts.LineStackFlows

	WangsuLiveLineStack      *echarts.LineStack
	WangsuLiveLineStackFlows *echarts.LineStackFlows

	WangsuLivePie      *echarts.Pie
	WangsuLivePieFlows *echarts.Pie

	WangsuCdnLineStack      *echarts.LineStack
	WangsuCdnLineStackFlows *echarts.LineStackFlows

	WangsuCdnPie      *echarts.Pie
	WangsuCdnPieFlows *echarts.Pie
}

func (v *WebValue) QWInit(conf *conf.Conf) error {
	// 饼图
	webValuePie := func(pie, pieflows *echarts.Pie, sql string, unit string) {
		pieflows.Sql = sql
		pieflows.Begin = v.Begin
		pieflows.End = v.End
		pieflows.Series = make([]echarts.PieSerie, 0)

		pieflows.Read(conf.Mysql)
		tmp, _ := json.Marshal(pieflows)

		_ = json.Unmarshal(tmp, pie)
		pie.SerieNameRatio(unit)
	}

	// 折线图
	webValueLineStack := func(linetackflows *echarts.LineStackFlows, sql string) *echarts.LineStack {
		linetackflows.Sql = sql
		linetackflows.Begin = v.Begin
		linetackflows.End = v.End
		linetackflows.Read(conf.Mysql)
		linetackflows.SumFlow()
		return linetackflows.ConvertLineStack()
	}

	// 七牛直播折线图
	v.QiniuLiveLineStack = new(echarts.LineStack)
	v.QiniuLiveLineStackFlows = new(echarts.LineStackFlows)
	v.QiniuLiveLineStack = webValueLineStack(
		v.QiniuLiveLineStackFlows,
		"SELECT hub,JSON_OBJECTAGG(date, JSON_EXTRACT(updown,'$.max')) FROM QiniuHubsFlow WHERE date >= ? AND date < ? GROUP BY hub HAVING SUM(JSON_EXTRACT(updown,'$.bytesum'))/POWER(1000,4)>0.1",
	)

	// 七牛直播饼图
	v.QiniuLivePie = new(echarts.Pie)
	v.QiniuLivePieFlows = new(echarts.Pie)
	webValuePie(
		v.QiniuLivePie,
		v.QiniuLivePieFlows,
		"SELECT hub,ROUND(SUM(JSON_EXTRACT(updown,'$.bytesum'))/POWER(1000,4),2) AS sumbyte FROM QiniuHubsFlow WHERE date >= ? AND date < ? GROUP BY hub HAVING sumbyte >0.1",
		"TB",
	)

	// 七牛 CDN 折线图
	v.QiniuCdnLineStack = new(echarts.LineStack)
	v.QiniuCdnLineStackFlows = new(echarts.LineStackFlows)
	v.QiniuCdnLineStack = webValueLineStack(
		v.QiniuCdnLineStackFlows,
		"SELECT domain,JSON_OBJECTAGG(date,ROUND(bandwidthmax/1000,0)) FROM QiniuCdnsFlow WHERE date >= ? AND date < ? GROUP BY domain HAVING SUM(bytesum)/POWER(1000,3) >0.1",
	)

	// 七牛CDN 饼图
	v.QiniuCdnPie = new(echarts.Pie)
	v.QiniuCdnPieFlows = new(echarts.Pie)
	webValuePie(
		v.QiniuCdnPie,
		v.QiniuCdnPieFlows,
		"SELECT domain,ROUND(SUM(bytesum)/POWER(1000,3),2) AS sum FROM QiniuCdnsFlow WHERE date >= ? AND date < ? GROUP BY domain HAVING sum >0.1",
		"GB",
	)

	// 网宿直播折线图
	v.WangsuLiveLineStack = new(echarts.LineStack)
	v.WangsuLiveLineStackFlows = new(echarts.LineStackFlows)
	v.WangsuLiveLineStack = webValueLineStack(
		v.WangsuLiveLineStackFlows,
		"SELECT channel,JSON_OBJECTAGG(date,peakValue) FROM WangsuLiveFlow WHERE date >= ? AND date < ? AND totalFlow/POWER(1024,2) > 0.1 GROUP BY channel",
	)

	// 网宿直播饼图
	v.WangsuLivePie = new(echarts.Pie)
	v.WangsuLivePieFlows = new(echarts.Pie)
	webValuePie(
		v.WangsuLivePie,
		v.WangsuLivePieFlows,
		"SELECT channel,SUM(totalFlow/POWER(1024,2)) AS total FROM WangsuLiveFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING total >0.1",
		"TB",
	)

	// 网宿 Cdn 折线图
	v.WangsuCdnLineStack = new(echarts.LineStack)
	v.WangsuCdnLineStackFlows = new(echarts.LineStackFlows)
	v.WangsuCdnLineStack = webValueLineStack(
		v.WangsuCdnLineStackFlows,
		"SELECT channel,JSON_OBJECTAGG(date,peakValue) FROM WangsuCdnFlow WHERE date >= ? AND date < ? AND totalFlow > 0.1 GROUP BY channel",
	)

	// cdn 饼图
	v.WangsuCdnPie = new(echarts.Pie)
	v.WangsuCdnPieFlows = new(echarts.Pie)
	webValuePie(
		v.WangsuCdnPie,
		v.WangsuCdnPieFlows,
		"SELECT channel,SUM(totalFlow) AS total FROM WangsuCdnFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING total >0.1",
		"GB",
	)
	
	return nil
}

// web 页面相关
func Start(conf *conf.Conf) {
	// 传递给 web 页面数据
	webValue := new(WebValue)
	webValue.End = time.Now()
	webValue.Begin = webValue.End.AddDate(0, 0, -3)

	webValue.QWInit(conf)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"json": func(s interface{}) string {
			jsonBytes, err := json.Marshal(s)
			if err != nil {
				return ""
			}
			return string(jsonBytes)
		},
	})
	r.Static("/template", "template/")
	r.LoadHTMLGlob("template/*.html")
	r.GET("/live", func(ctx *gin.Context) {
		ctx.HTML(200, "live.html", webValue)
	})
	r.GET("/cdn", func(ctx *gin.Context) {
		ctx.HTML(200, "cdn.html", webValue)
	})
	r.Run(":8173")
}
