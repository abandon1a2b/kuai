package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/abandon1a2b/kuai/util"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "git:stats",
		Short: "统计指定目录下所有 Git 项目的提交记录",
		Run:   runGitstatistic,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("path", "./", "work path")
	cmd.Flags().String("user", "user", "user")
	appendCommand(cmd)
}

const (
	lastYear = 1 // 统计过去一年的数据
)

func runGitstatistic(cmd *cobra.Command, _ []string) {
	root, _ := cmd.Flags().GetString("path")       // 指定根目录
	authorName, _ := cmd.Flags().GetString("user") // 指定根目录
	root, _ = util.AbsPath(root)

	var mu sync.Mutex
	err := util.WalkGitRepos(root, 10, func(projectPath string) {
		// 调用 Git 命令获取作者的提交次数和新增代码行数
		commitCount, addedLines, err := getAuthorCommitData(projectPath, authorName, lastYear)
		if err != nil {
			mu.Lock()
			fmt.Printf("Failed to get Git data for %s: %v\n", projectPath, err)
			mu.Unlock()
			return
		}

		// 统计作者在早上、下午和晚上提交的次数
		morningCount, afternoonCount, eveningCount, err := getAuthorCommitTimeData(projectPath, authorName, lastYear)
		if err != nil {
			mu.Lock()
			fmt.Printf("Failed to get Git time data for %s: %v\n", projectPath, err)
			mu.Unlock()
			return
		}

		mu.Lock()
		defer mu.Unlock()
		fmt.Println("-------------------------------------------------")
		fmt.Println("Git project found:", projectPath)
		fmt.Printf("Commit count: %d\n", commitCount)
		fmt.Printf("Added lines: %d\n", addedLines)
		totalCount := morningCount + afternoonCount + eveningCount
		if totalCount > 0 {
			fmt.Printf("Morning commits: %d %.2f%%\n", morningCount, float64(morningCount)*100.0/float64(totalCount))
			fmt.Printf("Afternoon commits: %d %.2f%%\n", afternoonCount, float64(afternoonCount)*100.0/float64(totalCount))
			fmt.Printf("Evening commits: %d %.2f%%\n", eveningCount, float64(eveningCount)*100.0/float64(totalCount))
		}
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
		return
	}
}

// 获取指定 Git 项目中作者的提交次数和新增代码行数
func getAuthorCommitData(projectPath, authorName string, year int) (int, int, error) {
	end := time.Now()
	start := end.AddDate(-year, 0, 0)

	cmd := exec.Command("git", "log", "--author="+authorName, "--since=\""+start.Format("2006-01-02")+"\"", "--before=\""+end.Format("2006-01-02")+"\"", "--shortstat")
	cmd.Dir = projectPath

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	commitCount := 0
	addedLines := 0

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			commitCount++
		} else if strings.Contains(line, "insertion") {
			parts := strings.Split(line, ",")
			if len(parts) > 0 {
				numStr := strings.TrimSpace(parts[1])
				num, _ := strconv.Atoi(numStr)
				addedLines += num
			}
		}
	}

	return commitCount, addedLines, nil
}

// 获取指定 Git 项目中作者在早上、下午和晚上提交的次数
func getAuthorCommitTimeData(projectPath, authorName string, year int) (int, int, int, error) {
	end := time.Now()
	start := end.AddDate(-year, 0, 0)

	cmd := exec.Command("git", "log", "--author="+authorName, "--since=\""+start.Format("2006-01-02")+"\"", "--before=\""+end.Format("2006-01-02")+"\"", "--pretty=format:%H|%cI")
	cmd.Dir = projectPath

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}

	morningCount := 0
	afternoonCount := 0
	eveningCount := 0

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}
		//commitHash := parts[0]
		commitTimeStr := parts[1]
		commitTime, err := time.Parse(time.RFC3339, commitTimeStr)
		if err != nil {
			continue
		}
		hour := commitTime.Hour()
		if hour >= 5 && hour < 12 { // 早上 5:00-11:59
			morningCount++
		} else if hour >= 12 && hour < 19 { // 下午 12:00-18:59
			afternoonCount++
		} else { // 晚上 19:00-4:59
			eveningCount++
		}
	}

	return morningCount, afternoonCount, eveningCount, nil
}
