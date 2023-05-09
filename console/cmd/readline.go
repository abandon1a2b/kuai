package cmd

import (
	"fmt"
	"github.com/leancodebox/goose/lineopt"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "tool:readline",
		Short: "文件读取",
		Run:   runReadline,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("path", "./readme.md", "path")
	appendCommand(cmd)
}

func runReadline(cmd *cobra.Command, _ []string) {
	filePath, _ := cmd.Flags().GetString("path")
	lineopt.ReadLine(filePath, func(item string) {
		fmt.Println(item)
	})
}
