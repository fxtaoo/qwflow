package main

import (
	"qwflow/timing"
	"qwflow/web"
)

func main() {
	// 定时执行，获取处理存储昨天流量
	go func() {
		timing.Start()
	}()

	web.Start()
}
