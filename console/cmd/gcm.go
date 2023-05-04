package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "tool:gcm",
		Short: "demo gcm",
		Run:   runGcm,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runGcm(_ *cobra.Command, _ []string) {
	t := time.Now()
	if t.Weekday() == time.Sunday || t.Weekday() == time.Saturday {
		fmt.Println("今天是周末")
	}
	timeData := mapToYesterday(t).Format(time.RFC3339)
	data := fmt.Sprintf(`export GIT_COMMITTER_DATE="%v" && export GIT_AUTHOR_DATE="%v"`, timeData, timeData)
	fmt.Println(data)
	fmt.Println("unset GIT_COMMITTER_DATE GIT_AUTHOR_DATE")
	/**
	打印一个时间，把今天的 00:00:00 ～ 24:00:00 映射到 前一天的 17:31:32 ~ 23:44:33, go 代码实现
	*/
}

func mapToYesterday(now time.Time) time.Time {
	yesterday := now.Add(-24 * time.Hour)

	// 计算昨天的开始时间和结束时间
	startYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 19, 31, 32, 0, yesterday.Location())
	endYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 44, 33, 0, yesterday.Location())

	// 计算映射时间在当天的位置
	diff := now.Sub(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
	ratio := float64(diff) / float64(24*time.Hour)

	// 映射到前一天的时间段
	return startYesterday.Add(time.Duration(int64(ratio*float64(endYesterday.Unix()-startYesterday.Unix()))) * time.Second)
}
