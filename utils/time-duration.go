package utils

import (
	"strconv"
	"time"
)

func GetHourDiffer(startTime, endTime string) int32 {
	var minutes int32
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		minutes = int32(diff / 60)
		return minutes
	} else {
		return minutes
	}
}

func GetTimeDiff(diff int32) string {
	if diff < 60 {
		return strconv.Itoa(int(diff)) + "分钟前"
	} else if diff > 1440 {
		return strconv.Itoa(int(diff/1440)) + "天前"
	} else {
		return strconv.Itoa(int(diff/60)) + "小时前"
	}
}
