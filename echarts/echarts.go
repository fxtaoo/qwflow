package echarts

import (
	"encoding/json"
	"fmt"
	"qwflow/mysql"
	"time"
)

// 要传给 web 页面 echarts 折线图数据
type LineStack struct {
	Title  string           `json:"title"`
	Legend []string         `json:"legend"`
	XAxis  []string         `json:"xAxis"`
	Series []LineStackSerie `json:"series"`
}

type LineStackSerie struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Stack string `json:"stack"`
	Data  []int  `json:"data"`
}

// 数据库折线图相关数据
type LineStackFlows struct {
	Sql   string
	Begin time.Time
	End   time.Time
	Flows []Flow
}

type Flow struct {
	Name        string         `json:"name"`
	DateFlowMax map[string]int `json:"dateflowmax"`
}

// 要传给 web 页面 echarts 饼图数据
type Pie struct {
	Sql    string     `json:"-"`
	Begin  time.Time  `json:"-"`
	End    time.Time  `json:"-"`
	Title  string     `json:"title"`
	Series []PieSerie `json:"series"`
}

type PieSerie struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
}

func (l *LineStackSerie) Init() {
	l.Type = "line"
	l.Stack = fmt.Sprint(time.Now().UnixNano())
	l.Data = make([]int, 0)
}

// Flows[] Name 添加前缀
func (l *LineStackFlows) SeriesNamePrefix(p string) {
	for i := range l.Flows {
		l.Flows[i].Name = p + "-" + l.Flows[i].Name
	}
}

// 添加一个汇总 Flow
func (l *LineStackFlows) SumFlow() {

	sumFlow := &Flow{Name: "汇总"}
	sumFlow.DateFlowMax = make(map[string]int)

	dateSlice := []string{}
	dateTmp := l.Begin
	for l.End.Sub(dateTmp).Hours()/24 > 0 {
		dateSlice = append(dateSlice, dateTmp.Format("2006-01-02"))
		dateTmp = dateTmp.AddDate(0, 0, 1)
	}

	for _, v := range dateSlice {
		for j := range l.Flows {
			sumFlow.DateFlowMax[v] += l.Flows[j].DateFlowMax[v]
		}
	}

	l.Flows = append(l.Flows, *sumFlow)

}

// 从数据库读数据
// 注意有表有固定的结构
func (l *LineStackFlows) Read(m mysql.Mysql) error {
	stmt, err := m.GetStmt(l.Sql)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(l.Begin.Format("20060102"),
		l.End.Format("20060102"))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		flowTmp := new(Flow)
		flowTmp.DateFlowMax = make(map[string]int)
		dateFlowMaxStr := ""

		err := rows.Scan(&flowTmp.Name, &dateFlowMaxStr)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(dateFlowMaxStr), &flowTmp.DateFlowMax)
		if err != nil {
			return err
		}
		l.Flows = append(l.Flows, *flowTmp)

	}
	return nil
}

// 数据 LineStackFlows 转换为 LineStack
func (l *LineStackFlows) ConvertLineStack() *LineStack {
	lineStack := new(LineStack)
	lineStack.Legend = make([]string, 0)
	lineStack.XAxis = make([]string, 0)
	lineStack.Series = make([]LineStackSerie, 0)

	dateTmp := l.Begin
	for l.End.Sub(dateTmp).Hours()/24 > 0 {
		lineStack.XAxis = append(lineStack.XAxis, dateTmp.Format("2006-01-02"))
		dateTmp = dateTmp.AddDate(0, 0, 1)
	}

	// 将 LineStackFlows 转换为 LineStack
	for i := range l.Flows {
		lineStack.Legend = append(lineStack.Legend, l.Flows[i].Name)
		lineStackSerie := new(LineStackSerie)
		lineStackSerie.Init()
		lineStackSerie.Name = l.Flows[i].Name
		for _, date := range lineStack.XAxis {
			lineStackSerie.Data = append(lineStackSerie.Data, l.Flows[i].DateFlowMax[date])
		}
		lineStack.Series = append(lineStack.Series, *lineStackSerie)
	}

	return lineStack
}

// 从数据库读数据
// 注意有表有固定的结构
func (p *Pie) Read(m mysql.Mysql) error {
	stmt, err := m.GetStmt(p.Sql)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(p.Begin.Format("20060102"),
		p.End.Format("20060102"))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		pieSerie := new(PieSerie)

		err := rows.Scan(&pieSerie.Name, &pieSerie.Value)
		if err != nil {
			return err
		}
		p.Series = append(p.Series, *pieSerie)
	}
	return nil
}

// 给命名名称加上值与百分比
func (p *Pie) SerieNameRatio(unit string) {
	var sum float64
	for _, v := range p.Series {
		sum += v.Value
	}

	for i := range p.Series {
		p.Series[i].Name = fmt.Sprintf("%s  %.2f %s  %.2f%%", p.Series[i].Name, p.Series[i].Value, unit, p.Series[i].Value/sum*100)
	}
}
