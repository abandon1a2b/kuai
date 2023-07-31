package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	cmd := &cobra.Command{
		Use:   "now",
		Short: "文件读取",
		Run:   getNow,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("path", "./readme.md", "path")
	appendCommand(cmd)
}

func getNow(cmd *cobra.Command, _ []string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
