package qiniu

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"qwflow/mysql"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/cdn"
	"github.com/qiniu/go-sdk/v7/pili"
)

// 直播相关
type Qiniu struct {
	Key struct {
		AccessKey string `json:"accesskey"`
		SecretKey string `json:"secretkey"`
	} `json:"key"`
	Mac  *auth.Credentials `json:"-"`
	Pili struct {
		Manager  *pili.Manager
		HubsFlow HubsFlow `json:"hubsflow"`
	} `json:"-"`
	Cdn struct {
		Manager  *cdn.CdnManager
		CndFlows *CndsFlows
	} `json:"-"`
}

type HubFlow struct {
	Name      string   `json:"name"`
	BeginDate string   `json:"begindate"`
	EndDate   string   `json:"enddate"`
	Up        FlowInfo `json:"up"`
	Down      FlowInfo `json:"down"`
	UpDown    FlowInfo `json:"updown"`
}

type FlowInfo struct {
	BandWidthMax         int    `json:"max"`
	BandWidthMaxDateTime string `json:"maxdatetime"`
	ByteSum              int    `json:"bytesum"`
}

type HubsFlow struct {
	Hubs []HubFlow `json:"hubs"`
}

// cdn 相关
type CndFlow struct {
	Domain       string `json:"domain"`
	Date         string `json:"date"`
	BandWidthMax int    `json:"bandwidthmax"`
	ByteSum      int    `json:"bytesum"`
}

type CndsFlows struct {
	Cnds []CndFlow
}

func (q *Qiniu) Init() error {
	// 直播相关
	q.Pili.Manager = pili.NewManager(pili.ManagerConfig{AccessKey: q.Key.AccessKey, SecretKey: q.Key.SecretKey})
	// 初始化 hub 列表
	hubs, err := q.Pili.Manager.GetHubList(context.Background())
	if err != nil {
		return err
	}
	for _, hub := range hubs.Items {
		q.Pili.HubsFlow.Hubs = append(q.Pili.HubsFlow.Hubs, HubFlow{Name: hub.Name})
	}

	// cdn 相关
	q.Mac = auth.New(q.Key.AccessKey, q.Key.SecretKey)
	// 初始化 cdn 域名列表
	err = q.CndDomainInit()
	if err != nil {
		return err
	}
	q.Cdn.Manager = cdn.NewCdnManager(q.Mac)
	return nil
}

// 初始化 cdn 列表
func (q *Qiniu) CndDomainInit() error {

	req, err := http.NewRequest("GET", "http://api.qiniu.com/domain?&limit=1000&SourceTypes=domain", nil)
	if err != nil {
		return err
	}

	accessToken, err := q.Mac.SignRequest(req)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "QBox "+accessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	domains := struct {
		Domains []struct {
			Name string `json:"name"`
		} `json:"domains"`
	}{}
	err = json.Unmarshal(resData, &domains)
	if err != nil {
		return err
	}
	q.Cdn.CndFlows = new(CndsFlows)
	q.Cdn.CndFlows.Cnds = make([]CndFlow, 0)
	for _, domain := range domains.Domains {
		q.Cdn.CndFlows.Cnds = append(q.Cdn.CndFlows.Cnds, CndFlow{Domain: domain.Name})
	}
	return nil
}

func (h *HubFlow) CalcFlow(manager *pili.Manager) error {
	if h.BeginDate == "" || h.EndDate == "" || manager == nil {
		return errors.New("h.BeginDate h.EndDate p.Manager 未全初始化")
	}
	GetStatCommonRequest := pili.GetStatCommonRequest{Begin: h.BeginDate, End: h.EndDate, G: "5min"}
	// 下行 Flow
	GetStatDownflowRequest := pili.GetStatDownflowRequest{GetStatCommonRequest: GetStatCommonRequest, Where: make(map[string][]string)}
	GetStatDownflowRequest.Where["hub"] = []string{h.Name}
	downFlow, error := manager.GetStatDownflow(context.Background(), GetStatDownflowRequest)
	if error != nil {
		return error
	}
	// 上行 Flow
	GetStatUpflowRequest := pili.GetStatUpflowRequest{GetStatCommonRequest: GetStatCommonRequest, Where: make(map[string][]string)}
	GetStatUpflowRequest.Where["hub"] = []string{h.Name}
	upFlow, error := manager.GetStatUpflow(context.Background(), GetStatUpflowRequest)
	if error != nil {
		return error
	}

	// 流量 byte
	for i := range downFlow {
		up := upFlow[i].Values["flow"]
		down := downFlow[i].Values["flow"]
		time := downFlow[i].Time.Format("2006-01-02 15:04:05")

		// hub 总流量
		h.Up.ByteSum += up
		h.Down.ByteSum += down

		// 最大上行流量
		if up > h.Up.BandWidthMax {
			h.Up.BandWidthMax = up
			h.Up.BandWidthMaxDateTime = time
		}

		// 最大下行流量
		if down > h.Down.BandWidthMax {
			h.Down.BandWidthMax = down
			h.Down.BandWidthMaxDateTime = time
		}

		// 最大总流量
		if up+down > h.UpDown.BandWidthMax {
			h.UpDown.BandWidthMax = up + down
			h.UpDown.BandWidthMaxDateTime = time
		}
	}
	h.Up.BandWidthMax = h.Up.BandWidthMax * 8 / (5 * 60) / 1000000
	h.Down.BandWidthMax = h.Down.BandWidthMax * 8 / (5 * 60) / 1000000
	h.UpDown.BandWidthMax = h.UpDown.BandWidthMax * 8 / (5 * 60) / 1000000
	h.UpDown.ByteSum = h.Up.ByteSum + h.Down.ByteSum

	return nil
}

// hub 数据保存到数据库
func (h *HubFlow) Save(m *mysql.Mysql, stmt *sql.Stmt) (int64, error) {

	up, _ := json.Marshal(h.Up)
	down, _ := json.Marshal(h.Down)
	updown, _ := json.Marshal(h.UpDown)

	return m.Write(stmt, h.Name, h.BeginDate, up, down, updown)
}

// 从七牛获取一天数据
func (h *HubsFlow) OneDayFlow(manager *pili.Manager, begin, end string) {

	var wg sync.WaitGroup
	for i := range h.Hubs {
		wg.Add(1)
		go func(i int) {
			h.Hubs[i].BeginDate = begin
			h.Hubs[i].EndDate = end
			// todo 错误处理
			h.Hubs[i].CalcFlow(manager)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// 保存数据到数据库
func (h *HubsFlow) Save(m *mysql.Mysql) error {
	// 查询配置，预处理
	stmt, err := m.GetStmt(`INSERT INTO QiniuHubsFlow(hub, date, up, down, updown) values (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var wg sync.WaitGroup
	for i := range h.Hubs {
		wg.Add(1)
		go func(i int) {
			// todo 错误处理
			h.Hubs[i].Save(m, stmt)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}

// 指定日期， 从七牛出抓取流量数据
// 日期格式 time.RFC3339 2022-10-08T00:00:00+08:00
func (h *HubsFlow) DataFlows(manager *pili.Manager, m *mysql.Mysql, begin, end time.Time) {
	// 前一天数据
	var wg sync.WaitGroup
	ch := make(chan struct{}, runtime.NumCPU())
	for end.Sub(begin).Hours()/24 > 0 {
		// 结构体深度拷贝
		hubsFlowTmp := new(HubsFlow)
		hubsFlowJson, _ := json.Marshal(h)
		_ = json.Unmarshal(hubsFlowJson, hubsFlowTmp)

		endTmp := begin.AddDate(0, 0, 1)

		// 限制并发数量
		ch <- struct{}{}
		wg.Add(1)
		go func(begin, end time.Time, h *HubsFlow) {
			h.OneDayFlow(manager, begin.Format("20060102"), endTmp.Format("20060102"))
			h.Save(m)
			wg.Done()
			<-ch
		}(begin, endTmp, hubsFlowTmp)
		begin = endTmp
	}
	wg.Wait()
}

// 获取流量
func (c *CndsFlows) DataFlows(manager *cdn.CdnManager, m *mysql.Mysql, begin, end time.Time) error {
	// 流量
	domains := []string{}

	for _, domain := range c.Cnds {
		domains = append(domains, domain.Domain)
	}
	end = end.AddDate(0, 0, -1)
	granularity := "day"
	// 流量
	fluxData, err := manager.GetFluxData(begin.Format("2006-01-02"), end.Format("2006-01-02"), granularity, domains)
	if err != nil {
		return err
	}
	// 带宽
	bandwidthData, err := manager.GetBandwidthData(begin.Format("2006-01-02"), end.Format("2006-01-02"), granularity, domains)
	if err != nil {
		return err
	}

	c.Cnds = make([]CndFlow, 0)

	// 日期
	for _, domain := range domains {
		for j, date := range fluxData.Time {
			tmpCndFlow := new(CndFlow)
			tmpCndFlow.Domain = domain
			tmpCndFlow.Date = strings.Split(date, " ")[0]
			// 流量
			if fluxData.Data[tmpCndFlow.Domain].DomainChina != nil {
				tmpCndFlow.ByteSum += fluxData.Data[tmpCndFlow.Domain].DomainChina[j]
			}
			if fluxData.Data[tmpCndFlow.Domain].DomainOversea != nil {
				tmpCndFlow.ByteSum += fluxData.Data[tmpCndFlow.Domain].DomainOversea[j]
			}
			// 带宽
			if bandwidthData.Data[tmpCndFlow.Domain].DomainChina != nil {
				tmpCndFlow.BandWidthMax += bandwidthData.Data[tmpCndFlow.Domain].DomainChina[j]
			}
			if bandwidthData.Data[tmpCndFlow.Domain].DomainOversea != nil {
				tmpCndFlow.BandWidthMax += bandwidthData.Data[tmpCndFlow.Domain].DomainOversea[j]
			}
			c.Cnds = append(c.Cnds, *tmpCndFlow)
		}
	}

	// 保存到数据库
	err = c.Save(m)
	if err != nil {
		return err
	}

	return nil
}

// 保存到数据库
func (c *CndsFlows) Save(m *mysql.Mysql) error {
	// 查询配置，预处理
	stmt, err := m.GetStmt(`INSERT INTO QiniuCdnsFlow(domain, date, bandwidthmax, bytesum) values (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	ch := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	for i := range c.Cnds {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			// todo 错误处理
			m.Write(stmt, c.Cnds[i].Domain, c.Cnds[i].Date, c.Cnds[i].BandWidthMax, c.Cnds[i].ByteSum)
			wg.Done()
			<-ch
		}(i)
	}
	wg.Wait()
	return nil
}
