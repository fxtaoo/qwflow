package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"qwflow/conf"
	"qwflow/echarts"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WebValue struct {
	Begin time.Time
	End   time.Time
	Name  string

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

	QiniuWangsuLiveLineStack *echarts.LineStack
	QiniuWangsuLivePie       *echarts.Pie

	QiniuWangsuCdnLineStack *echarts.LineStack
	QiniuWangsuCdnPie       *echarts.Pie
}

func (v *WebValue) QWLiveInit(conf *conf.Conf) error {
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
		"SELECT hub,ROUND(SUM(JSON_EXTRACT(updown,'$.bytesum'))/POWER(1000,4),2) AS sumbyte FROM QiniuHubsFlow WHERE date >= ? AND date < ? GROUP BY hub HAVING sumbyte >1",
		"TB",
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
		"SELECT channel,SUM(ROUND(totalFlow/POWER(1024,2),2)) AS total FROM WangsuLiveFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING total >0.1",
		"TB",
	)

	// 汇总折线图
	v.QiniuWangsuLiveLineStack = new(echarts.LineStack)
	v.QiniuLiveLineStackFlows.SeriesNamePrefix("七牛")
	v.WangsuLiveLineStackFlows.SeriesNamePrefix("网宿")
	v.QiniuLiveLineStackFlows.Flows = append(
		v.QiniuLiveLineStackFlows.Flows,
		v.WangsuLiveLineStackFlows.Flows...,
	)
	v.QiniuWangsuLiveLineStack = v.QiniuLiveLineStackFlows.ConvertLineStack()

	// 汇总饼图
	v.QiniuWangsuLivePie = new(echarts.Pie)
	v.QiniuLivePieFlows.SeriesNamePrefix("七牛")
	v.WangsuLivePieFlows.SeriesNamePrefix("网宿")
	v.QiniuWangsuLivePie.Series = append(
		v.QiniuWangsuLivePie.Series,
		v.QiniuLivePieFlows.Series...,
	)
	v.QiniuWangsuLivePie.Series = append(
		v.QiniuWangsuLivePie.Series,
		v.WangsuLivePieFlows.Series...,
	)
	v.QiniuWangsuLivePie.SerieNameRatio("TB")

	return nil

}

func (v *WebValue) QWCdnInit(conf *conf.Conf) error {
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

	// 七牛 CDN 折线图
	v.QiniuCdnLineStack = new(echarts.LineStack)
	v.QiniuCdnLineStackFlows = new(echarts.LineStackFlows)
	v.QiniuCdnLineStack = webValueLineStack(
		v.QiniuCdnLineStackFlows,
		"SELECT domain,JSON_OBJECTAGG(date,ROUND(bandwidthmax/1024/1024,0)) FROM QiniuCdnsFlow WHERE date >= ? AND date < ? GROUP BY domain HAVING SUM(bytesum)/POWER(1024,3) >1",
	)

	// 七牛CDN 饼图
	v.QiniuCdnPie = new(echarts.Pie)
	v.QiniuCdnPieFlows = new(echarts.Pie)
	webValuePie(
		v.QiniuCdnPie,
		v.QiniuCdnPieFlows,
		"SELECT domain,ROUND(SUM(bytesum)/POWER(1024,3),2) AS sum FROM QiniuCdnsFlow WHERE date >= ? AND date < ? GROUP BY domain HAVING sum >1",
		"GB",
	)

	// 网宿 Cdn 折线图
	v.WangsuCdnLineStack = new(echarts.LineStack)
	v.WangsuCdnLineStackFlows = new(echarts.LineStackFlows)
	v.WangsuCdnLineStack = webValueLineStack(
		v.WangsuCdnLineStackFlows,
		"SELECT channel,JSON_OBJECTAGG(date,peakValue) FROM WangsuCdnFlow WHERE date >= ? AND date < ? AND totalFlow > 1 GROUP BY channel",
	)

	// cdn 饼图
	v.WangsuCdnPie = new(echarts.Pie)
	v.WangsuCdnPieFlows = new(echarts.Pie)
	webValuePie(
		v.WangsuCdnPie,
		v.WangsuCdnPieFlows,
		"SELECT channel,SUM(totalFlow) AS total FROM WangsuCdnFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING total >1",
		"GB",
	)

	// 汇总折线图
	v.QiniuWangsuCdnLineStack = new(echarts.LineStack)
	v.QiniuCdnLineStackFlows.SeriesNamePrefix("七牛")
	v.WangsuCdnLineStackFlows.SeriesNamePrefix("网宿")
	v.QiniuCdnLineStackFlows.Flows = append(
		v.QiniuCdnLineStackFlows.Flows,
		v.WangsuCdnLineStackFlows.Flows...,
	)
	v.QiniuWangsuCdnLineStack = v.QiniuCdnLineStackFlows.ConvertLineStack()

	// 汇总饼图
	v.QiniuWangsuCdnPie = new(echarts.Pie)
	v.QiniuCdnPieFlows.SeriesNamePrefix("七牛")
	v.WangsuCdnPieFlows.SeriesNamePrefix("网宿")
	v.QiniuWangsuCdnPie.Series = append(
		v.QiniuWangsuCdnPie.Series,
		v.QiniuCdnPieFlows.Series...,
	)
	v.QiniuWangsuCdnPie.Series = append(
		v.QiniuWangsuCdnPie.Series,
		v.WangsuCdnPieFlows.Series...,
	)
	v.QiniuWangsuCdnPie.SerieNameRatio("GB")

	return nil
}

func (v *WebValue) DateSelect(ctx *gin.Context) {
	month := 1
	if ctx.Query("month") != "" {
		month, _ = strconv.Atoi(ctx.Query("month"))
		v.Name = fmt.Sprintf("%d 月", month)
	}

	day := 0
	if ctx.Query("day") != "" {
		day, _ = strconv.Atoi(ctx.Query("day"))
		month = 0
		v.Name = "1 周"
	}

	v.End = time.Now()
	v.Begin = v.End.AddDate(0, -month, -day)

	// 自定义日期选择
	if ctx.Query("begen") != "" {
		begen, _ := time.Parse("2006-01-02 -0700 MST", ctx.Query("begen")+" +0800 CST")
		end, _ := time.Parse("2006-01-02 -0700 MST", ctx.Query("end")+" +0800 CST")

		v.Begin = begen
		v.End = end.AddDate(0, 0, 1)

		v.Name = fmt.Sprintf(
			"%s/%s 总计 %d 天",
			ctx.Query("begen"),
			ctx.Query("end"),
			int(end.Sub(v.Begin).Hours()/24)+1,
		)
	}

}

// web 页面相关
func Start() {
	var conf conf.Conf
	// 初始化数据
	err := conf.Init()
	if err != nil {
		// todo log 处理
		log.Fatal(err)
	}
	defer conf.Mysql.DB.Close()

	r := gin.Default()

	r.Use(gin.BasicAuth(conf.Accounts))

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
		webValue := new(WebValue)
		webValue.DateSelect(ctx)

		webValue.QWLiveInit(&conf)

		ctx.HTML(200, "live.html", webValue)
	})

	r.GET("/cdn", func(ctx *gin.Context) {
		webValue := new(WebValue)
		webValue.DateSelect(ctx)

		webValue.QWCdnInit(&conf)

		ctx.HTML(200, "cdn.html", webValue)
	})
	r.Run(":8173")
}
