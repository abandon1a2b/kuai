package cmd

import (
	"fmt"
	"github.com/abandon1a2b/kuai/util"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	cmd := &cobra.Command{
		Use:   "git:scanRepo",
		Short: "git 指定目录下全项目拉取 git:scanRepo --path=./",
		Run:   runGitScanRepo,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("path", "./", "work path")
	appendCommand(cmd)
}

func runGitScanRepo(cmd *cobra.Command, _ []string) {
	root, _ := cmd.Flags().GetString("path") // 指定根目录
	root, _ = util.AbsPath(root)
	// 定义一个匿名函数，用于处理每个目录
	var visitDirFunc = func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() && path != root {
			fmt.Println("Subdirectory found:", path)
			gitRepo(path)
			return filepath.SkipDir // 只遍历当前目录下的子目录，不递归遍历子目录下的子目录
		}
		return nil
	}

	// 递归遍历指定目录下的所有文件和子目录
	err := filepath.Walk(root, visitDirFunc)
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
		return
	}
}

func gitRepo(path string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 判断是否为 Git 项目
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		fmt.Println("Git project found:", path)
		// 调用 git pull 命令
		cmd := exec.Command("git", "-C", path, "remote", "-v")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return
		}
		fmt.Println(string(output))
	}
}
