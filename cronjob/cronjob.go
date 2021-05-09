package cronjob

import (
	"GoProject/fudan_bbs/dal"
	"time"
)

func InitCronJob() {
	// 启动定时任务
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			dal.LoadThreadsToRedis()
		}
	}()
}
