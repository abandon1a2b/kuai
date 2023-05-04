package cmd

import (
	"bufio"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/leancodebox/goose/fileopt"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "tool:pic2",
		Short: "点图生成",
		Run:   runPicTool2,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runPicTool2(_ *cobra.Command, _ []string) {

	fileopt.DirExistOrCreate("storage/pic/")
	dc := gg.NewContext(17747, 997)
	f, errF := os.OpenFile("/Users/thh/Downloads/go_frame/output/result.txt", os.O_RDONLY, 0666)
	if errF != nil {
		fmt.Print(errF)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	counter := 1
	for scanner.Scan() {
		if counter == 1 {
			counter += 1
			continue
		}
		line := scanner.Text()
		fmt.Println(line)
		allT := strings.Split(string(line), "(")
		fmt.Println(allT)
		for _, item := range allT[1:] {
			fmt.Println(item)
			str := strings.ReplaceAll(item, ")", "")
			xy := strings.Split(str, ",")
			dc.LineTo(cast.ToFloat64(xy[0]), cast.ToFloat64(xy[1]))
		}

		dc.SetColor(randColor())
		dc.SetFillRule(gg.FillRuleEvenOdd)
		dc.FillPreserve()
		dc.SetColor(randColor())
		dc.SetLineWidth(1)
		dc.Stroke()
		fmt.Println("end1")
	}
	err := dc.SavePNG("out.png")
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println("finish")
	}
}
