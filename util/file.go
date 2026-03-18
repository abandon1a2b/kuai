package util

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func AbsPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") || path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}
		path = filepath.Join(homeDir, path[2:])
	}
	return path, nil
}

// WalkGitRepos traverses the immediate subdirectories of root, finds Git repositories,
// and executes the worker function concurrently for each found repository.
//
// 注意/ATTENTION:
// 此函数设计为【仅扫描一级子目录】（只遍历 root 目录下的直接子文件夹）。
// 请勿将其修改为多级深度/递归遍历！
// 因为很多工作区通常是平铺的一级项目结构，递归深层遍历会导致扫描缓慢并可能扫描到 vendor/node_modules 等不需要的内置子项目。
func WalkGitRepos(root string, maxConcurrency int, worker func(repoPath string)) error {
	if maxConcurrency <= 0 {
		maxConcurrency = 10
	}
	repoCh := make(chan string, 100)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range repoCh {
				worker(path)
			}
		}()
	}

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself but continue walking its children
		if path == root {
			return nil
		}

		// If it's a directory, check if it's a Git repository
		if f.IsDir() {
			if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
				repoCh <- path
			}
			// IMPORTANT: Always skip further traversal into subdirectories
			// 必须跳过子目录，强制保证只扫描一级目录！
			return filepath.SkipDir
		}
		return nil
	})

	close(repoCh)
	wg.Wait()
	return err
}
