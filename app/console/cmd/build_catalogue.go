package cmd

import (
	"bytes"
	"fmt"
	"github.com/leancodebox/goose/fileopt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "tool:build_catalogue",
		Short: "生成文件目录",
		Run:   runBuildCatalogue,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runBuildCatalogue(_ *cobra.Command, _ []string) {
	path := `config.GetString("BUILD_MD_LIST")`

	list := ScanPathBuildList(path)
	bgd(list, ".", 0)
	err := fileopt.FilePutContents(path+"/"+"readme.new.md", listBuffer.String())
	if err != nil {
		fmt.Println(err)
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
			t.Type = 1
			t.Name = file.Name()
			t.Children = ScanPathBuildList(dirName + string(os.PathSeparator) + file.Name())
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
func ScanDir(dirName string) []string {
	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Println(err)
	}
	var fileList []string
	for _, file := range files {
		fileList = append(fileList, dirName+string(os.PathSeparator)+file.Name())
	}
	return fileList
}

// ScanDirs 递归扫描目录
func ScanDirs(dirName string) []string {
	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Println(err)
	}
	var fileList []string
	for _, file := range files {
		fileList = append(fileList, dirName+string(os.PathSeparator)+file.Name())
		if file.IsDir() {
			fileList = append(fileList, ScanDir(dirName+string(os.PathSeparator)+file.Name())...)
		}
	}
	return fileList
}
