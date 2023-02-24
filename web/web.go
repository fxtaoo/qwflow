package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"qwflow/conf"
	"qwflow/echarts"
	"qwflow/timing"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WebValue struct {
	Begin time.Time
	End   time.Time
	Name  string

	CdnOtherGB  int
	LiveOtherTB float64
	DownloadImg bool

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
	webValuePie := func(pie, pieflows *echarts.Pie, sql string, unit string) string {
		pieflows.Sql = sql
		pieflows.Begin = v.Begin
		pieflows.End = v.End
		pieflows.Series = make([]echarts.PieSerie, 0)

		pieflows.Read(conf.Mysql)
		// (p *Pie) AddOther(otherGB int)
		otherNames := pieflows.AddOther(int(v.LiveOtherTB * 1000))
		tmp, _ := json.Marshal(pieflows)

		_ = json.Unmarshal(tmp, pie)
		pie.SerieNameRatio(unit)

		return otherNames
	}

	// 折线图
	webValueLineStack := func(linetackflows *echarts.LineStackFlows, sql string, otherNames string, pie *echarts.Pie) *echarts.LineStack {
		linetackflows.Sql = sql
		linetackflows.Begin = v.Begin
		linetackflows.End = v.End
		linetackflows.Read(conf.Mysql)
		linetackflows.AddOther(otherNames, pie)
		linetackflows.SumFlow()
		return linetackflows.ConvertLineStack()
	}
	// 七牛直播饼图
	v.QiniuLivePie = new(echarts.Pie)
	v.QiniuLivePieFlows = new(echarts.Pie)
	otherNames := webValuePie(
		v.QiniuLivePie,
		v.QiniuLivePieFlows,
		"SELECT hub,ROUND(SUM(JSON_EXTRACT(updown,'$.bytesum'))/POWER(1000,4),2) AS sumbyte FROM QiniuHubsFlow WHERE date >= ? AND date < ? GROUP BY hub HAVING sumbyte > 0.1 ORDER BY sumbyte DESC",
		"TB",
	)

	// 七牛直播折线图
	v.QiniuLiveLineStack = new(echarts.LineStack)
	v.QiniuLiveLineStackFlows = new(echarts.LineStackFlows)
	v.QiniuLiveLineStack = webValueLineStack(
		v.QiniuLiveLineStackFlows,
		"SELECT hub,JSON_OBJECTAGG(date, JSON_EXTRACT(updown,'$.max')),AVG(JSON_EXTRACT(updown,'$.max')) AS avg FROM QiniuHubsFlow WHERE date >= ? AND date < ? GROUP BY hub HAVING avg > 0.1 ORDER BY avg DESC",
		otherNames,
		v.QiniuLivePie,
	)

	// 网宿直播饼图
	v.WangsuLivePie = new(echarts.Pie)
	v.WangsuLivePieFlows = new(echarts.Pie)
	otherNames = webValuePie(
		v.WangsuLivePie,
		v.WangsuLivePieFlows,
		"SELECT channel,SUM(ROUND(totalFlow/POWER(1024,2),2)) AS total FROM WangsuLiveFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING total >0.1 ORDER BY total DESC",
		"TB",
	)

	// 网宿直播折线图
	v.WangsuLiveLineStack = new(echarts.LineStack)
	v.WangsuLiveLineStackFlows = new(echarts.LineStackFlows)
	v.WangsuLiveLineStack = webValueLineStack(
		v.WangsuLiveLineStackFlows,
		"SELECT channel,JSON_OBJECTAGG(date,peakValue), AVG(peakValue) AS avg FROM WangsuLiveFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING avg > 1 ORDER BY avg DESC",
		otherNames,
		v.WangsuLivePie,
	)

	// 汇总折线图
	v.QiniuWangsuLiveLineStack = new(echarts.LineStack)
	v.QiniuLiveLineStackFlows.SeriesNamePrefix("七牛")
	v.WangsuLiveLineStackFlows.SeriesNamePrefix("网宿")

	tmp := new(echarts.LineStackFlows)
	tmp.Begin = v.Begin
	tmp.End = v.End
	tmp.OrderAdd(v.QiniuLiveLineStackFlows, v.WangsuLiveLineStackFlows)
	v.QiniuWangsuLiveLineStack = tmp.ConvertLineStack()

	// 汇总饼图
	v.QiniuWangsuLivePie = new(echarts.Pie)
	v.QiniuLivePieFlows.SeriesNamePrefix("七牛")
	v.WangsuLivePieFlows.SeriesNamePrefix("网宿")
	v.QiniuWangsuLivePie.OrderAdd(v.QiniuLivePieFlows, v.WangsuLivePieFlows)
	v.QiniuWangsuLivePie.SerieNameRatio("TB")

	return nil

}

func (v *WebValue) QWCdnInit(conf *conf.Conf) error {
	// 饼图
	webValuePie := func(pie, pieflows *echarts.Pie, sql string, unit string) string {
		pieflows.Sql = sql
		pieflows.Begin = v.Begin
		pieflows.End = v.End
		pieflows.Series = make([]echarts.PieSerie, 0)

		pieflows.Read(conf.Mysql)
		otherNames := pieflows.AddOther(v.CdnOtherGB)
		tmp, _ := json.Marshal(pieflows)

		_ = json.Unmarshal(tmp, pie)
		pie.SerieNameRatio(unit)

		return otherNames
	}

	// 折线图
	webValueLineStack := func(linetackflows *echarts.LineStackFlows, sql string, otherNames string, pie *echarts.Pie) *echarts.LineStack {
		linetackflows.Sql = sql
		linetackflows.Begin = v.Begin
		linetackflows.End = v.End
		linetackflows.Read(conf.Mysql)
		linetackflows.AddOther(otherNames, pie)
		linetackflows.SumFlow()
		return linetackflows.ConvertLineStack()
	}

	// 七牛CDN 饼图
	v.QiniuCdnPie = new(echarts.Pie)
	v.QiniuCdnPieFlows = new(echarts.Pie)
	otherNames := webValuePie(
		v.QiniuCdnPie,
		v.QiniuCdnPieFlows,
		"SELECT domain,ROUND(SUM(bytesum)/POWER(1024,4),5) AS sum FROM QiniuCdnsFlow WHERE date >= ? AND date < ? GROUP BY domain HAVING sum >0.001 ORDER BY sum DESC",
		"TB",
	)

	// 七牛 CDN 折线图
	v.QiniuCdnLineStack = new(echarts.LineStack)
	v.QiniuCdnLineStackFlows = new(echarts.LineStackFlows)
	v.QiniuCdnLineStack = webValueLineStack(
		v.QiniuCdnLineStackFlows,
		"SELECT domain,JSON_OBJECTAGG(date,ROUND(bandwidthmax/1024/1024,0)),AVG(ROUND(bandwidthmax/1024/1024,0)) AS avg FROM QiniuCdnsFlow WHERE date >= ? AND date < ? GROUP BY domain HAVING avg > 1 ORDER BY avg DESC",
		otherNames,
		v.QiniuCdnPie,
	)

	// 网宿 cdn 饼图
	v.WangsuCdnPie = new(echarts.Pie)
	v.WangsuCdnPieFlows = new(echarts.Pie)
	otherNames = webValuePie(
		v.WangsuCdnPie,
		v.WangsuCdnPieFlows,
		"SELECT channel,SUM(totalFlow)/1024 AS total FROM WangsuCdnFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING total >0.001 ORDER BY total DESC",
		"TB",
	)

	// 网宿 Cdn 折线图
	v.WangsuCdnLineStack = new(echarts.LineStack)
	v.WangsuCdnLineStackFlows = new(echarts.LineStackFlows)
	v.WangsuCdnLineStack = webValueLineStack(
		v.WangsuCdnLineStackFlows,
		"SELECT channel,JSON_OBJECTAGG(date,peakValue),AVG(peakValue) AS avg FROM WangsuCdnFlow WHERE date >= ? AND date < ? GROUP BY channel HAVING avg >1 ORDER BY avg DESC",
		otherNames,
		v.WangsuCdnPie,
	)

	// 汇总饼图
	v.QiniuWangsuCdnPie = new(echarts.Pie)
	v.QiniuCdnPieFlows.SeriesNamePrefix("七牛")
	v.WangsuCdnPieFlows.SeriesNamePrefix("网宿")
	v.QiniuWangsuCdnPie.OrderAdd(v.QiniuCdnPieFlows, v.WangsuCdnPieFlows)
	v.QiniuWangsuCdnPie.SerieNameRatio("TB")

	// 汇总折线图
	v.QiniuWangsuCdnLineStack = new(echarts.LineStack)
	v.QiniuCdnLineStackFlows.SeriesNamePrefix("七牛")
	v.WangsuCdnLineStackFlows.SeriesNamePrefix("网宿")

	tmp := new(echarts.LineStackFlows)
	tmp.Begin = v.Begin
	tmp.End = v.End
	tmp.OrderAdd(v.QiniuCdnLineStackFlows, v.WangsuCdnLineStackFlows)
	v.QiniuWangsuCdnLineStack = tmp.ConvertLineStack()

	return nil
}

func (v *WebValue) DateSelect(ctx *gin.Context) {
	v.Name = "1 个月"
	month := 1
	if ctx.Query("month") != "" {
		month, _ = strconv.Atoi(ctx.Query("month"))
		v.Name = fmt.Sprintf("%d 个月", month)
	}

	day := 0
	if ctx.Query("day") != "" {
		day, _ = strconv.Atoi(ctx.Query("day"))
		month = 0
		if day%7 == 0 {
			if day/7 == 2 {
				v.Name = "半月"
			} else {
				v.Name = fmt.Sprintf("%d 周", day/7)
			}
		} else {
			v.Name = fmt.Sprintf("%d 天", day)
		}

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

	// cdn 聚合筛选值
	if ctx.Query("cdnOtherGB") != "" {
		v.CdnOtherGB, _ = strconv.Atoi(ctx.Query("cdnOtherGB"))
	}

	// live 聚合筛选值
	if ctx.Query("liveOtherTB") != "" {
		v.LiveOtherTB, _ = strconv.ParseFloat(ctx.Query("liveOtherTB"), 64)
	}

	// 是否下载图片
	if ctx.Query("downloadImg") == "true" {
		v.DownloadImg = true
	}
}

// 再一次获取数据
func (v *WebValue) getDataAgainResult(ctx *gin.Context) error {
	var conf conf.Conf
	// 初始化数据
	err := conf.Init()
	if err != nil {
		return err
	}
	// 数据库初始化
	conf.Mysql.Init()
	defer conf.Mysql.DB.Close()

	bengin, _ := time.Parse("2006-01-02 -0700 MST", ctx.Query("date")+" +0800 CST")
	end := bengin.AddDate(0, 0, 1)
	err = timing.GetDayFlow(&conf, bengin, end, ctx.Query("sort"), ctx.Request.URL.Path[1:], true)
	if err != nil {
		return err
	}

	return nil
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

	getDataAgainResult := make(chan string)

	getDataAgain := func(webValue *WebValue, ctx *gin.Context) {
		if ctx.Query("sort") != "" {
			quit := make(chan int)
			go func() {
				// get updateresult 失败避免锁死
				for {
					select {
					case <-quit:
						return
					case <-time.After(time.Second * 3):
						<-getDataAgainResult
					}
				}

			}()
			err := webValue.getDataAgainResult(ctx)
			if err != nil {
				getDataAgainResult <- fmt.Sprintf("获取失败：%s", err.Error())
			} else {
				getDataAgainResult <- "获取成功"
			}
			quit <- 0
			// 延迟避免不出现弹窗
			<-time.After(time.Millisecond * 1500)
		}
	}

	r.GET("/live", func(ctx *gin.Context) {
		// 数据库初始化
		conf.Mysql.Init()
		defer conf.Mysql.DB.Close()

		webValue := new(WebValue)
		webValue.LiveOtherTB = conf.LiveOtherTB
		getDataAgain(webValue, ctx)
		webValue.DateSelect(ctx)
		webValue.QWLiveInit(&conf)

		ctx.HTML(200, "live.html", webValue)
	})

	r.GET("/cdn", func(ctx *gin.Context) {
		// 数据库初始化
		conf.Mysql.Init()
		defer conf.Mysql.DB.Close()

		webValue := new(WebValue)
		webValue.CdnOtherGB = conf.CdnOtherGB
		getDataAgain(webValue, ctx)
		webValue.DateSelect(ctx)

		webValue.QWCdnInit(&conf)

		ctx.HTML(200, "cdn.html", webValue)
	})

	r.GET("/getDataAgainResult", func(ctx *gin.Context) {
		quit := make(chan int)
		go func() {
			// get updateresult 失败避免锁死
			for {
				select {
				case <-quit:
					return
				case <-time.After(time.Second * 6):
					getDataAgainResult <- "获取失败"
				}
			}

		}()
		ctx.String(200, "%s", <-getDataAgainResult)
		quit <- 0
	})
	r.Run(":8174")
}
