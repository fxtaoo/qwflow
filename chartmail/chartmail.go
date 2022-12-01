package chartmail

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"qwflow/mail"
	"strings"
)

// 流量报表
type ChartMail struct {
	Switch     bool            `json:"switch"`
	Stmp       mail.Smtp       `json:"smtp"`
	Mail       []string        `json:"mail"`
	BodyEnd    string          `json:"bodyEnd"`
	Body       strings.Builder `json:"-"`
	ImgDirPath string          `json:"imgDirPath"`
}

func (c *ChartMail) SendMail() []error {

	m := mail.Mail{
		Subject: "七牛网宿直播 cdn 流量带宽报表",
		To:      c.Mail,
	}

	sort := []string{"live", "cdn"}
	dm := []string{"d7", "m1", "m3"}

	for _, v1 := range sort {
		c.Body.WriteString(fmt.Sprintf("<br><h2>%s 相关</h2><br>", v1))
		for _, v2 := range dm {
			c.BodyAddImg(fmt.Sprintf("%s-%s-qiniu-wangsu-line-stack.png", v1, v2))
			c.BodyAddImg(fmt.Sprintf("%s-%s-qiniu-wangsu-pie.png", v1, v2))
		}
	}

	// 说明放至末尾
	c.Body.WriteString(c.BodyEnd)

	m.Body = c.Body.String()

	err := m.SendAlone(&c.Stmp)
	if err != nil {
		return err
	}

	return nil
}

// 图片转 base64
func (c *ChartMail) BodyAddImg(fileName string) {
	filePath := path.Join(c.ImgDirPath, fileName)
	fileByte, _ := os.ReadFile(filePath)
	src := base64.StdEncoding.EncodeToString(fileByte)
	c.Body.WriteString("<img src=\"data:image/png;base64,")
	c.Body.WriteString(src)
	c.Body.WriteString("\" />")
}
