package util

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"log"
)

// FormatGfTime 格式化时间
func FormatGfTime(gt *gtime.Time) string {
	if gt == nil {
		return ""
	}
	now := gtime.Now().Timestamp()
	timestamp := gt.Timestamp()
	diff := now - timestamp

	const (
		yearInSeconds   = 31536000
		dayInSeconds    = 86400
		hourInSeconds   = 3600
		minuteInSeconds = 60
		secondInSeconds = 1
	)

	log.Printf("进入了===%v", diff)
	switch {
	case diff > yearInSeconds:
		return fmt.Sprintf("%d年前", int(diff/yearInSeconds))
	case diff > dayInSeconds:
		return fmt.Sprintf("%d天前", int(diff/dayInSeconds))
	case diff > hourInSeconds:
		return fmt.Sprintf("%d小时前", int(diff/hourInSeconds))
	case diff > minuteInSeconds:
		return fmt.Sprintf("%d分钟前", int(diff/minuteInSeconds))
	case diff > secondInSeconds:
		return fmt.Sprintf("%d秒前", int(diff/secondInSeconds))
	default:
		return "刚刚"
	}
}
