package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
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
	timeData := ""
	if t.Weekday() != time.Sunday && t.Weekday() != time.Saturday {
		minute := cast.ToString(60 / 24 * cast.ToInt(t.Format("15")))
		timeData = t.Format("2006-01-02T") + "07:" + minute + t.Format(":05Z07:00")
	} else {
		timeData = t.Format(time.RFC3339)
	}
	data := fmt.Sprintf(`export GIT_COMMITTER_DATE="%v" && export GIT_AUTHOR_DATE="%v"`, timeData, timeData)
	fmt.Println(data)
	fmt.Println("unset GIT_COMMITTER_DATE GIT_AUTHOR_DATE")
}
