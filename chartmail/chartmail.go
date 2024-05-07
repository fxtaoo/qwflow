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
	ImgName    string          `json:"imgName"`
}

func (c *ChartMail) SendMail(sort, imgName string) []error {

	m := mail.Mail{
		Subject: fmt.Sprintf("七牛网宿 %s 带宽流量报表", sort),
		To:      c.Mail,
	}

	c.Body.WriteString(fmt.Sprintf("<br><h2>%s 相关</h2><br>", sort))

	c.BodyAddImg(fmt.Sprintf("%s-%s-qiniu-wangsu-line-stack.png", sort, imgName))
	c.BodyAddImg(fmt.Sprintf("%s-%s-qiniu-wangsu-pie.png", sort, imgName))

	c.BodyAddImg(fmt.Sprintf("%s-%s-qiniu-line-stack.png", sort, imgName))
	c.BodyAddImg(fmt.Sprintf("%s-%s-qiniu-pie.png", sort, imgName))

	c.BodyAddImg(fmt.Sprintf("%s-%s-wangsu-line-stack.png", sort, imgName))
	c.BodyAddImg(fmt.Sprintf("%s-%s-wangsu-pie.png", sort, imgName))

	// 说明放至末尾
	c.Body.WriteString(c.BodyEnd)
	m.Body = c.Body.String()
	defer c.Body.Reset()

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
