// 流量日环比增幅超过设定值邮件告警
package alert

import (
	"fmt"
	"qwflow/mail"
	"qwflow/mysql"
	"time"
)

// 一天警告警告
type Alerts struct {
	CdnGrowthWarnNum  float64 `json:"cdnGrowthWarnNum"`
	LiveGrowthWarnNum float64 `json:"liveGrowthWarnNum"`
	Date              time.Time
	Alerts            []Alert
	Mail              []string  `json:"mail"`
	Stmp              mail.Smtp `json:"smtp"`
}

// 一条警告
type Alert struct {
	Name       string
	Sort       string
	Unit       string
	YYesterday float64
	Yesterday  float64
	Growth     float64
}

func (a *Alerts) Init() {
	now := time.Now()
	a.Date = now
	a.Alerts = make([]Alert, 0)
}

// 从数据库读数据，并比较
func (a *Alerts) Calc(m *mysql.Mysql) error {

	yyesterday := a.Date.AddDate(0, 0, -2).Format("2006-01-02")
	yesterday := a.Date.AddDate(0, 0, -1).Format("2006-01-02")

	calcfunc := func(a *Alerts, growthWarnNum float64, sort, unit, sql string) error {
		stmt, err := m.GetStmt(sql)
		if err != nil {
			return err
		}
		rows, err := stmt.Query(
			yyesterday,
			yesterday,
			growthWarnNum,
		)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmpAlert := Alert{
				Sort: sort,
				Unit: unit,
			}
			err := rows.Scan(&tmpAlert.Name, &tmpAlert.YYesterday, &tmpAlert.Yesterday, &tmpAlert.Growth)
			if err != nil {
				return err
			}

			a.Alerts = append(a.Alerts, tmpAlert)
		}

		return nil
	}

	// 七牛直播
	err := calcfunc(
		a,
		a.LiveGrowthWarnNum,
		"七牛直播",
		"TB",
		"SELECT tmp.hub,tmp.n1,tmp.n2,ROUND(tmp.n2/tmp.n1,2) AS growth FROM ( SELECT q1.hub, ROUND(JSON_EXTRACT(q1.updown,'$.bytesum')/POWER(1000,4),2) AS n1, ROUND(JSON_EXTRACT(q2.updown,'$.bytesum')/POWER(1000,4),2) AS n2 FROM QiniuHubsFlow q1,QiniuHubsFlow q2 WHERE q1.hub=q2.hub AND q1.date=? AND q2.date=? ) tmp HAVING tmp.n2>0.1 AND growth >?",
	)
	if err != nil {
		return err
	}

	// 七牛 cdn
	err = calcfunc(
		a,
		a.CdnGrowthWarnNum,
		"七牛 cdn",
		"GB",
		"SELECT tmp.domain,tmp.n1,tmp.n2,ROUND(tmp.n2/tmp.n1,2) AS growth FROM ( SELECT q1.domain, ROUND(q1.bytesum/POWER(1024,3),2) AS n1, ROUND(q2.bytesum/POWER(1024,3),2) AS n2 FROM QiniuCdnsFlow q1,QiniuCdnsFlow q2 WHERE q1.domain=q2.domain AND q1.date=? AND q2.date=? ) tmp HAVING tmp.n2 >1 AND growth >?",
	)
	if err != nil {
		return err
	}

	// 网宿直播
	err = calcfunc(
		a,
		a.LiveGrowthWarnNum,
		"网宿直播",
		"TB",
		"SELECT tmp.channel,tmp.n1,tmp.n2,ROUND(tmp.n2/tmp.n1,2) AS growth FROM ( SELECT w1.channel, ROUND(w1.totalFlow/POWER(1000,2),2) AS n1, ROUND(w2.totalFlow/POWER(1000,2),2) AS n2 FROM WangsuLiveFlow w1,WangsuLiveFlow w2 WHERE w1.channel=w2.channel AND w1.date=? AND w2.date=? ) tmp HAVING tmp.n2>0.1 AND growth >?",
	)
	if err != nil {
		return err
	}

	// 网宿 cdn
	err = calcfunc(
		a,
		a.CdnGrowthWarnNum,
		"网宿 cdn",
		"GB",
		"SELECT w1.channel,w1.totalFlow,w2.totalFlow, ROUND(w2.totalFlow/w1.totalFlow,2) AS growth FROM WangsuCdnFlow w1,WangsuCdnFlow w2 WHERE w1.channel=w2.channel AND w1.date=? AND w2.date=? AND w2.totalFlow>1 HAVING growth >?",
	)
	if err != nil {
		return err
	}
	return nil
}

// 发送邮件
func (a *Alerts) SendMail() {
	// 聚合信息
	subject := "七牛网宿直播 cdn 流量日增幅告警！"
	content := ""

	yyesterday := a.Date.AddDate(0, 0, -2).Format("2006-01-02")
	yesterday := a.Date.AddDate(0, 0, -1).Format("2006-01-02")

	for _, v := range a.Alerts {
		content += fmt.Sprintf("<br>%s-%s (%s %.2f %s) (%s %.2f %s) 增长 %d%%",
			v.Sort,
			v.Name,
			yyesterday,
			v.YYesterday,
			v.Unit,
			yesterday,
			v.Yesterday,
			v.Unit,
			int(v.Growth*100),
		)
	}
	if content != "" {
		// 发送邮件
		m := mail.Mail{
			Subject: subject,
			Body:    content,
			To:      a.Mail,
		}
		m.Send(&a.Stmp)
	}
}
