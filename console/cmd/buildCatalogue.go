package cmd

import (
	"bytes"
	"fmt"
	"github.com/eatmeatball/kuai/util"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "tool:build_catalogue",
		Short: "生成文件目录",
		Run:   runBuildCatalogue,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("path", "./", "dir path")
	cmd.Flags().String("output", "", "output")
	appendCommand(cmd)
}

func runBuildCatalogue(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("path") // 指定根目录
	path, _ = util.AbsPath(path)
	output, _ := cmd.Flags().GetString("output") // 指定根目录

	list := ScanPathBuildList(path)
	bgd(list, ".", 0)
	if output != "" {
		err := os.WriteFile(output, listBuffer.Bytes(), 0644)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(fmt.Sprintf("output %s success\n", output))
		}
	} else {
		fmt.Println(listBuffer.String())
	}
}

var listBuffer = bytes.Buffer{}

func bgd(list []PNode, prefix string, deep int) {
	for _, node := range list {
		listBuffer.WriteString("\n")
		if deep != 0 {
			listBuffer.WriteString(StrPad(" ", deep*2))
		}
		//M::$gd .= '- [' . trimBlank($key) . '](' . trimBlank($prefix  . '/' . $key, '%20') . ')';

		rList := []string{" ", "　", "\n", "\r", "\t"}
		itemName := StrReplaces(rList, "_", node.Name)
		pathName := StrReplaces(rList, "%20", prefix+string(os.PathSeparator)+node.Name)

		listBuffer.WriteString(fmt.Sprintf("- [%v](%v)", itemName, pathName))
		if node.Type == 1 {
			bgd(node.Children, prefix+string(os.PathSeparator)+node.Name, deep+1)
		}
	}
}

// StrReplaces 多字段替换
func StrReplaces(sList []string, n string, itemName string) string {
	for _, s := range sList {
		itemName = strings.ReplaceAll(itemName, s, n)
	}
	return itemName
}

// StrPad 字符串补充
func StrPad(padString string, padLength int) string {
	padBuffer := bytes.Buffer{}
	for i := 0; i < padLength; i++ {
		padBuffer.WriteString(padString)
	}
	return padBuffer.String()
}

type PNode struct {
	Type     int
	Name     string
	Children []PNode
}

func ScanPathBuildList(dirName string) []PNode {
	var node []PNode
	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Println(err)
		return node
	}
	for _, file := range files {
		t := PNode{}
		isIgnore := false
		for _, item := range []string{"..", ".", ".git", ".idea"} {
			if file.Name() == item {
				isIgnore = true
				break
			}
		}
		if isIgnore {
			continue
		}
		if file.IsDir() {
			// 如果子目录没有可以生成目录的文件，那么直接跳过
			subNode := ScanPathBuildList(dirName + string(os.PathSeparator) + file.Name())
			if len(subNode) == 0 {
				continue
			}
			t.Type = 1
			t.Name = file.Name()
			t.Children = subNode
		} else {
			if !strings.Contains(strings.ToLower(file.Name()), ".md") {
				continue
			}
			t.Type = 0
			t.Name = file.Name()
		}
		node = append(node, t)
	}
	return node
}
