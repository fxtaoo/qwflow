package wangsu

import (
	"encoding/json"
	"fmt"
	"qwflow/mysql"
	"qwflow/wangsu/api/client"
	"qwflow/wangsu/common/auth"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WangSu struct {
	Key struct {
		AccessKey string `json:"accesskey"`
		SecretKey string `json:"secretkey"`
	} `json:"key"`
	AkskConfig    auth.AkskConfig `json:"-"`
	ReChannelName []struct {
		Suffix  string   `json:"suffix"`
		Channel []string `json:"channel"`
	} `json:"reChannelName"`
}

// 单条数据
type ChannelPeak struct {
	Channel   string
	Date      string
	PeakValue float64
	TotalFlow float64
}

// 一段日期数据
type DateChannelPeak struct {
	Begin            time.Time
	End              time.Time
	ChannelPeakSlice []ChannelPeak
}

type LiveData struct {
	Domains   []string
	TotalFlow float64
	PeakValue float64
}

// 初始化
func (w *WangSu) Init() {
	w.AkskConfig.AccessKey = w.Key.AccessKey
	w.AkskConfig.SecretKey = w.Key.SecretKey
	w.AkskConfig.Method = "POST"
}

// 按日期获取流量
func (d *DateChannelPeak) DataChannelPeak(w *WangSu, accetype string) error {
	w.AkskConfig.Uri = "/myview/bandwidth-peak-ranking"
	// 间隔天数
	dayNum := d.End.Sub(d.Begin).Hours() / 24
	chChannelPeakSlice := make(chan []ChannelPeak)

	// 控制并发数
	ch := make(chan struct{}, runtime.NumCPU())

	for i := 0; i < int(dayNum); i++ {
		ch <- struct{}{}
		go func(i int) {

			request := client.BandwidthPeakRankingRequest{}
			request.SetDataformat("json")
			request.SetAccetype(accetype)
			date := d.Begin.AddDate(0, 0, i).Format("2006-01-02")
			request.SetStartdate(date)
			request.SetEnddate(date)

			response := auth.Invoke(w.AkskConfig, request.String())

			// 网宿返回数据解析
			var vesponseValue struct {
				Provider struct {
					Date struct {
						ChannelPeak []struct {
							Channel   string `json:"channel"`
							PeakTime  string `json:"peakTime"`
							PeakValue string `json:"peakValue"`
							TotalFlow string `json:"totalFlow"`
						} `json:"channelPeak"`
					} `json:"date"`
				} `json:"provider"`
			}

			_ = json.Unmarshal(response, &vesponseValue)
			tmpChChannelPeakSlice := make([]ChannelPeak, 0)
			for _, v := range vesponseValue.Provider.Date.ChannelPeak {
				tmpChannelPeak := &ChannelPeak{
					Channel: v.Channel,
				}
				tmpChannelPeak.Date = strings.Split(v.PeakTime, " ")[0]
				tmpChannelPeak.PeakValue, _ = strconv.ParseFloat(v.PeakValue, 64)
				tmpChannelPeak.TotalFlow, _ = strconv.ParseFloat(v.TotalFlow, 64)
				// 剔除零值
				if tmpChannelPeak.TotalFlow == 0 {
					continue
				}
				tmpChChannelPeakSlice = append(tmpChChannelPeakSlice, *tmpChannelPeak)
			}
			chChannelPeakSlice <- tmpChChannelPeakSlice
			<-ch
		}(i)
	}

	for i := 0; i < int(dayNum); i++ {
		d.ChannelPeakSlice = append(d.ChannelPeakSlice, <-chChannelPeakSlice...)
	}
	return nil
}

// CDN 相关流量可直接用 DataChannelPeak
// 直播相关流量需要聚合几个域名同一时刻，最大带宽和，需要不同的接口
func (d *DateChannelPeak) LiveDataChannelPeak(w *WangSu, m *mysql.Mysql) error {

	// 讲同属一业务的域名聚合
	tmpD := DateChannelPeak{
		Begin: d.Begin,
		End:   d.Begin.AddDate(0, 0, 1),
	}
	// /myview/bandwidth-peak-ranking 接口获取
	tmpD.DataChannelPeak(w, "livestream")

	livesData := make(map[string]*LiveData)

	for i := range tmpD.ChannelPeakSlice {
		name := tmpD.ChannelPeakSlice[i].Channel
		// 只保留主域名
		// 同一业务相同主域名，多个不通二级推拉流域名
		name = name[strings.Index(name, ".")+1:]

		if _, ok := livesData[name]; !ok {
			livesData[name] = &LiveData{}
		}

		livesData[name].Domains = append(livesData[name].Domains, tmpD.ChannelPeakSlice[i].Channel)
	}

	dayNum := d.End.Sub(d.Begin).Hours() / 24
	for i := 0; i < int(dayNum); i++ {
		date := d.Begin.AddDate(0, 0, i)
		for key, v := range livesData {
			v.CalcPeakValue(w, key, date)
		}
		for key, v := range livesData {
			d.ChannelPeakSlice = append(d.ChannelPeakSlice, ChannelPeak{
				Channel:   key,
				Date:      date.Format("2006-01-02"),
				TotalFlow: v.TotalFlow,
				PeakValue: v.PeakValue,
			})
			// 置零
			v.TotalFlow = 0
		}
	}

	// 存储到数据库
	d.Save(m, "WangsuLiveFlow")
	return nil

}

// 计算单个业务多域名带宽峰值
// 时间为 date 当天
func (l *LiveData) CalcPeakValue(w *WangSu, name string, date time.Time) error {
	w.AkskConfig.Uri = "/api/report/bandwidth/multi-domain/real-time/edge"

	// todo 域名是不是有重复
	domains := []*string{}
	for i := range l.Domains {
		domains = append(domains, &l.Domains[i])
	}

	dateTime, _ := time.Parse(time.RFC3339, date.Format("2006-01-02")+"T00:00:00+08:00")

	chTotal := make(chan float64)
	chMaxValue := make(chan float64)
	ch := make(chan struct{}, runtime.NumCPU())

	// 一天 24 小时
	for i := 0; i < 24; i++ {
		ch <- struct{}{}
		go func(i int, dateTime time.Time, domains []*string) {
			tmpMaxValue := 0.0
			tmpTotal := 0.0
			request := client.ReportBandwidthMultiDomainRealTimeEdgeServiceRequest{}
			request.SetDataInterval("5m")
			request.SetDomain(domains)
			request.SetDateFrom(dateTime.Format(time.RFC3339))
			request.SetDateTo(dateTime.Add(time.Hour).Format(time.RFC3339))

			response := auth.Invoke(w.AkskConfig, request.String())

			var responseDate struct {
				Data []struct {
					Total         float64 `json:"total,string"`
					BandwidthData []struct {
						Value float64 `json:"value,string"`
					} `json:"bandwidthData"`
				} `json:"data"`
			}
			_ = json.Unmarshal(response, &responseDate)

			// 寻最大值
			if len(responseDate.Data) > 0 {
				for _, v := range responseDate.Data[0].BandwidthData {
					if v.Value > tmpMaxValue {
						tmpMaxValue = v.Value
					}
				}
				tmpTotal = responseDate.Data[0].Total
			}
			<-ch
			chTotal <- tmpTotal
			chMaxValue <- tmpMaxValue

		}(i, dateTime, domains)
		dateTime = dateTime.Add(time.Hour)
	}

	maxValue := 0.0
	for i := 0; i < 24; i++ {
		l.TotalFlow += <-chTotal
		tmp := <-chMaxValue
		if tmp > maxValue {
			maxValue = tmp
		}

	}

	l.PeakValue = maxValue
	return nil
}

// 常规是一个直播业务，一个主域名，多个二级域名
// 从数据库拉取信息，通过聚合相同主域名，将相同主域名视为同一业务
// 存在特殊情况，两个业务同一个主域名，这种情况需要标记主域名
func (d *DateChannelPeak) ReChannelName(w *WangSu) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, runtime.NumCPU())
	for i := range w.ReChannelName {
		wg.Add(1)
		ch <- struct{}{}
		go func(i int) {
			// 所有标记域名聚合成字符串
			tmpStr := strings.Join(w.ReChannelName[i].Channel, " ")
			for j := range d.ChannelPeakSlice {
				// 属于聚合字符串字串追加后缀标记
				if strings.Contains(tmpStr, d.ChannelPeakSlice[j].Channel) {
					d.ChannelPeakSlice[j].Channel += "-" + w.ReChannelName[i].Suffix
				}
			}
			wg.Done()
			<-ch
		}(i)
	}
	wg.Wait()
}

// 保存到数据库
func (d *DateChannelPeak) Save(m *mysql.Mysql, table string) error {
	// 查询配置，预处理
	sql := fmt.Sprintf("INSERT INTO %s(channel, date, peakValue, totalFlow) values (?, ?, ?, ?)", table)
	stmt, err := m.GetStmt(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ch := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	for i := range d.ChannelPeakSlice {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			m.Write(
				stmt,
				d.ChannelPeakSlice[i].Channel,
				d.ChannelPeakSlice[i].Date,
				d.ChannelPeakSlice[i].PeakValue,
				d.ChannelPeakSlice[i].TotalFlow,
			)
			wg.Done()
			<-ch
		}(i)
	}
	wg.Wait()
	return nil
}
