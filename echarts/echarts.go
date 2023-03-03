package echarts

import (
	"encoding/json"
	"fmt"
	"qwflow/mysql"
	"strings"
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
	Name      string `json:"name"`
	Type      string `json:"type"`
	Stack     string `json:"stack"`
	Data      []int  `json:"data"`
	MarkPoint struct {
		Data []MarkPointData `json:"data"`
	} `json:"markPoint"`
	MarkLine struct {
		Data      []MarkLineData `json:"data"`
		LineStyle struct {
			Width int    `json:"width"`
			Type  string `json:"type"`
		} `json:"lineStyle"`
	} `json:"markLine"`
	Smooth     bool `json:"smooth"`
	ShowSymbol bool `json:"showSymbol"`
}

type MarkPointData struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	SymbolRotate int    `json:"symbolRotate"`
	ItemStyle    struct {
		Color string `json:"color"`
	} `json:"itemStyle"`
}

type MarkLineData struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Label struct {
		Formatter  string `json:"formatter"`
		Position   string `json:"position"`
		FontWeight string `json:"fontWeight"`
	} `json:"label"`
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
	Avg         float64        `json:"-"`
	SumTag      bool           // 集合值标注，true 为集合值
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

func (l *LineStackSerie) Init(stack string) {
	l.Type = "line"
	l.Stack = stack
	l.Data = make([]int, 0)
	l.Smooth = true
	l.ShowSymbol = false
}

// 汇总数据添加标注
func (l *LineStackSerie) AddLabel() {
	// 最大最小
	l.MarkPoint.Data = []MarkPointData{
		{
			Type: "max",
			Name: "最大值",
			ItemStyle: struct {
				Color string "json:\"color\""
			}{
				Color: "rgba(244, 80, 80, 1)",
			},
		},
		{
			Type: "min",
			Name: "最小值",
			ItemStyle: struct {
				Color string "json:\"color\""
			}{
				Color: "rgba(38, 135, 37, 1)",
			},
		},
	}

	// 平均值
	l.MarkLine.Data = []MarkLineData{
		{
			Type: "average",
			Name: "平均值",
			Label: struct {
				Formatter  string "json:\"formatter\""
				Position   string "json:\"position\""
				FontWeight string `json:"fontWeight"`
			}{
				Formatter:  fmt.Sprintf("%s 平均值 {c} Mbps", l.Name),
				Position:   "middle",
				FontWeight: "bold",
			},
		},
	}

	l.MarkLine.LineStyle.Width = 2
	l.MarkLine.LineStyle.Type = "dashed"
}

// Flows[] Name 添加前缀
func (l *LineStackFlows) SeriesNamePrefix(p string) {
	for i := range l.Flows {
		l.Flows[i].Name = p + "-" + l.Flows[i].Name
	}
}

// 将较小的聚合成一个其他
func (l *LineStackFlows) AddOther(otherNames string, pie *Pie) {

	// 饼图，折线图数据筛选尺度不一致，导致饼图没有，折线图有
	// 以饼图为准
	// 数据以排序过，相比饼图多出数据并入其他
	lenLine := len(l.Flows)
	num := lenLine - len(pie.Series)
	for i := 1; i <= num; i++ {
		otherNames = fmt.Sprintf("%s %s", otherNames, l.Flows[lenLine-i].Name)
	}

	newFlows := make([]Flow, 0)
	otherFlows := make([]Flow, 0)

	for _, v := range l.Flows {
		if strings.Contains(otherNames, v.Name) {
			otherFlows = append(otherFlows, v)
		} else {
			newFlows = append(newFlows, v)
		}
	}

	l.Flows = otherFlows

	l.SumFlow()

	otherFlow := l.Flows[0]
	otherFlow.Name = "其他"
	otherFlow.SumTag = false

	// 其他插入合适位置
	i := 0
	for ; i < len(newFlows); i++ {
		if otherFlow.Avg > newFlows[i].Avg {
			break
		}
	}

	il := len(newFlows)
	l.Flows = make([]Flow, 0)

	if i == il {
		l.Flows = append(newFlows, otherFlow)
	} else {
		l.Flows = append(l.Flows, newFlows[0:i]...)
		l.Flows = append(l.Flows, otherFlow)
		l.Flows = append(l.Flows, newFlows[i:il]...)
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
			sumFlow.Avg += l.Flows[j].Avg
		}
	}

	sumFlow.Avg = sumFlow.Avg / float64(len(sumFlow.DateFlowMax))

	sumFlow.SumTag = true

	tmp := make([]Flow, 0)
	tmp = append(tmp, *sumFlow)
	tmp = append(tmp, l.Flows...)
	l.Flows = tmp
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

		err := rows.Scan(&flowTmp.Name, &dateFlowMaxStr, &flowTmp.Avg)
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
		lineStackSerie.Init(fmt.Sprint(i))
		lineStackSerie.Name = l.Flows[i].Name

		// 是否为集合值
		if l.Flows[i].SumTag {
			lineStackSerie.AddLabel()
		}

		for _, date := range lineStack.XAxis {
			lineStackSerie.Data = append(lineStackSerie.Data, l.Flows[i].DateFlowMax[date])
		}
		lineStack.Series = append(lineStack.Series, *lineStackSerie)
	}

	return lineStack
}

// 有序添加，降序
// l1 l2 已经是降序
func (l *LineStackFlows) OrderAdd(l1, l2 *LineStackFlows) {
	l1l := len(l1.Flows)
	l2l := len(l2.Flows)
	ll := l1l + l2l

	l1i := 0
	l2i := 0

	for i := 0; i < ll; i++ {
		if l1i == l1l {
			l.Flows = append(l.Flows, l2.Flows[l2i])
			l2i++
			continue
		}

		if l2i == l2l {
			l.Flows = append(l.Flows, l1.Flows[l1i])
			l1i++
			continue
		}

		if l1.Flows[l1i].Avg > l2.Flows[l2i].Avg {
			l.Flows = append(l.Flows, l1.Flows[l1i])
			l1i++
		} else {
			l.Flows = append(l.Flows, l2.Flows[l2i])
			l2i++
		}
	}

}

// Flows[] Name 添加前缀
func (p *Pie) SeriesNamePrefix(pstr string) {
	for i := range p.Series {
		p.Series[i].Name = pstr + "-" + p.Series[i].Name
	}
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
		p.Series[i].Name = fmt.Sprintf("%s %.1f%s %.1f%%", p.Series[i].Name, p.Series[i].Value, unit, p.Series[i].Value/sum*100)
	}
}

func (p *Pie) AddOther(otherGB int) string {
	// 小于这个值聚合为其他
	min := p.End.Sub(p.Begin).Hours() / 24 / 1000 * float64(otherGB)
	otherNames := ""
	newSeries := make([]PieSerie, 0)
	otherPieSerie := PieSerie{Name: "其他"}

	for _, e := range p.Series {
		if e.Value < min {
			otherNames += " " + e.Name
			otherPieSerie.Value += e.Value
		} else {
			newSeries = append(newSeries, e)
		}
	}

	// newSeries 已经是降序
	i := 0
	l := len(newSeries)

	for ; i < l; i++ {
		if otherPieSerie.Value > newSeries[i].Value {
			break
		}
	}

	p.Series = make([]PieSerie, 0)
	if i == l {
		p.Series = append(newSeries, otherPieSerie)
	} else {
		p.Series = append(p.Series, newSeries[0:i]...)
		p.Series = append(p.Series, otherPieSerie)
		p.Series = append(p.Series, newSeries[i:l]...)
	}

	return otherNames
}

// 有序添加，降序
// p1 p2 已经是降序
func (p *Pie) OrderAdd(p1, p2 *Pie) {

	p1l := len(p1.Series)
	p2l := len(p2.Series)

	pl := p1l + p2l

	p1i := 0
	p2i := 0

	for i := 0; i < pl; i++ {
		if p1i == p1l {
			p.Series = append(p.Series, p2.Series[p2i])
			p2i++
			continue
		}

		if p2i == p2l {
			p.Series = append(p.Series, p1.Series[p1i])
			p1i++
			continue
		}

		if p1.Series[p1i].Value > p2.Series[p2i].Value {
			p.Series = append(p.Series, p1.Series[p1i])
			p1i++
		} else {
			p.Series = append(p.Series, p2.Series[p2i])
			p2i++
		}
	}
}
