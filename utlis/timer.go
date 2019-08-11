package utlis

import (
	"time"
)

func Runtime(interval time.Duration, callback func()) {
	//创建定时器并设置定时时间
	TimerDemo := time.NewTimer(interval* time.Second)

	//循环监听定时器
	for {
		select {
		case <-TimerDemo.C:
			callback()
			//超时后重置定时器
			TimerDemo.Reset(interval * time.Second)
		}
	}
}


