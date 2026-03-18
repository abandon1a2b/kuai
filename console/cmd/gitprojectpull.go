package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/abandon1a2b/kuai/util"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "git:allpull",
		Short: "批量拉取指定目录下的所有 Git 项目",
		Run:   runAllprojectpull,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("path", "./", "work path")
	appendCommand(cmd)
}

func runAllprojectpull(cmd *cobra.Command, _ []string) {
	root, _ := cmd.Flags().GetString("path") // 指定根目录
	root, _ = util.AbsPath(root)

	err := util.WalkGitRepos(root, 10, func(path string) {
		gitPull(path)
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
		return
	}
}

func gitPull(path string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 判断是否为 Git 项目
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		fmt.Println("Git project found:", path)
		// 调用 git pull 命令
		cmd := exec.Command("git", "-C", path, "pull")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return
		}
		fmt.Println(string(output))
	}
}
