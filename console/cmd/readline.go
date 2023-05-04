package cmd

import (
	"fmt"
	"github.com/leancodebox/goose/lineopt"
	"github.com/spf13/cobra"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "tool:readline",
		Short: "文件读取",
		Run:   runReadline,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runReadline(_ *cobra.Command, _ []string) {
	lineopt.ReadLine("./readme.md", func(item string) {
		fmt.Println(item)
	})
}
