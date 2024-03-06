package cmd

import (
	"fmt"
	"github.com/eatmeatball/kuai/util"
	"github.com/leancodebox/goose/fileopt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "tool:pic",
		Short: "小图片生成",
		Run:   runPicTool,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().Int("h", 512, "height -h=512")
	cmd.Flags().Int("w", 256, "width -w=512")
	cmd.Flags().String("o", "./storage/pic/", "output -o=./storage/pic")
	cmd.Flags().Int("n", 10, "number -n=10")
	appendCommand(cmd)
}

func runPicTool(cmd *cobra.Command, _ []string) {

	h, _ := cmd.Flags().GetInt("h")
	w, _ := cmd.Flags().GetInt("w")
	o, _ := cmd.Flags().GetString("o")
	n, _ := cmd.Flags().GetInt("n")

	var S = h
	var S2 = w
	var maxT = 3

	outputRoot, _ := util.AbsPath(o)
	outputRoot = ensurePathHasTrailingSlash(outputRoot)
	fileopt.DirExistOrCreate(outputRoot)

	dc := createWatermark(S2)
	err := dc.SavePNG(outputRoot + "/water.png")
	if err != nil {
		fmt.Println(err)
	}

	water, _ := gg.LoadImage("storage/pic/water.png")
	x := h // 图片长(px)
	y := h // 图片宽(px)

	// 水印位置：以原图长宽 - 去水印图 - 偏移像素
	waterPosition := image.Pt(
		(x)/2,
		(y)/2,
	)

	for i := 1; i <= n; i++ {
		dc = gg.NewContext(x, y)
		dc.DrawImage(water, waterPosition.X-32, waterPosition.Y-32)
		dc.SetColor(randColor())
		dc.DrawCircle(float64(S/maxT), float64(S/maxT), float64(S/maxT))
		dc.Fill()
		dc.SetColor(randColor())
		dc.DrawString(time.Now().Format("2006_01_02_15_04_05")+cast.ToString(i), float64(S/maxT), float64(S/maxT))
		dc.Fill()

		filename := "out" + cast.ToString(time.Now().Format("2006_01_02_15_04_05")) + cast.ToString(i) + ".png"
		filePath := outputRoot + filename
		err = dc.SavePNG(filePath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("success :", filePath)
		}
	}
}

func ensurePathHasTrailingSlash(path string) string {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	return path
}

// 创建水印图
func createWatermark(size int) *gg.Context {
	dc := gg.NewContext(size, size)

	dc.DrawEllipse(float64(size/2), float64(size/3), float64(size/2), float64(size/2))
	dc.SetRGB(255, 255, 255)
	dc.SetColor(color.Black)
	dc.Fill()

	return dc
}

func randColor() *color.NRGBA {
	colorLen := len(colorList()) - 1
	return colorList()[rand.Intn(colorLen)]
}

func colorList() []*color.NRGBA {
	return []*color.NRGBA{Red50, Red100, Red200, Red300, Red400, Red500, Red600, Red700, Red800, Red900, Pink50, Pink100, Pink200, Pink300, Pink400, Pink500, Pink600, Pink700, Pink800, Pink900, Purple50, Purple100, Purple200, Purple300, Purple400, Purple500, Purple600, Purple700, Purple800, Purple900, Indigo50, Indigo100, Indigo200, Indigo300, Indigo400, Indigo500, Indigo600, Indigo700, Indigo800, Indigo900, Blue50, Blue100, Blue200, Blue300, Blue400, Blue500, Blue600, Blue700, Blue800, Blue900, Cyan50, Cyan100, Cyan200, Cyan300, Cyan400, Cyan500, Cyan600, Cyan700, Cyan800, Cyan900, Teal50, Teal100, Teal200, Teal300, Teal400, Teal500, Teal600, Teal700, Teal800, Teal900, Green50, Green100, Green200, Green300, Green400, Green500, Green600, Green700, Green800, Green900, Yellow50, Yellow100, Yellow200, Yellow300, Yellow400, Yellow500, Yellow600, Yellow700, Yellow800, Yellow900, Orange50, Orange100, Orange200, Orange300, Orange400, Orange500, Orange600, Orange700, Orange800, Orange900, Brown50, Brown100, Brown200, Brown300, Brown400, Brown500, Brown600, Brown700, Brown800, Brown900, Gray50, Gray100, Gray200, Gray300, Gray400, Gray500, Gray600, Gray700, Gray800, Gray900, BlueGray50, BlueGray100, BlueGray200, BlueGray300, BlueGray400, BlueGray500, BlueGray600, BlueGray700, BlueGray800, BlueGray900}
}

var Red50 = &color.NRGBA{R: 0xff, G: 0xeb, B: 0xee, A: 0xff}
var Red100 = &color.NRGBA{R: 0xff, G: 0xcd, B: 0xd2, A: 0xff}
var Red200 = &color.NRGBA{R: 0xef, G: 0x9a, B: 0x9a, A: 0xff}
var Red300 = &color.NRGBA{R: 0xe5, G: 0x73, B: 0x73, A: 0xff}
var Red400 = &color.NRGBA{R: 0xef, G: 0x53, B: 0x50, A: 0xff}
var Red500 = &color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
var Red600 = &color.NRGBA{R: 0xe5, G: 0x39, B: 0x35, A: 0xff}
var Red700 = &color.NRGBA{R: 0xd3, G: 0x2f, B: 0x2f, A: 0xff}
var Red800 = &color.NRGBA{R: 0xc6, G: 0x28, B: 0x28, A: 0xff}
var Red900 = &color.NRGBA{R: 0xb7, G: 0x1c, B: 0x1c, A: 0xff}
var Pink50 = &color.NRGBA{R: 0xfc, G: 0xe4, B: 0xec, A: 0xff}
var Pink100 = &color.NRGBA{R: 0xf8, G: 0xbb, B: 0xd0, A: 0xff}
var Pink200 = &color.NRGBA{R: 0xf4, G: 0x8f, B: 0xb1, A: 0xff}
var Pink300 = &color.NRGBA{R: 0xf0, G: 0x62, B: 0x92, A: 0xff}
var Pink400 = &color.NRGBA{R: 0xec, G: 0x40, B: 0x7a, A: 0xff}
var Pink500 = &color.NRGBA{R: 0xe9, G: 0x1e, B: 0x63, A: 0xff}
var Pink600 = &color.NRGBA{R: 0xd8, G: 0x1b, B: 0x60, A: 0xff}
var Pink700 = &color.NRGBA{R: 0xc2, G: 0x18, B: 0x5b, A: 0xff}
var Pink800 = &color.NRGBA{R: 0xad, G: 0x14, B: 0x57, A: 0xff}
var Pink900 = &color.NRGBA{R: 0x88, G: 0x0e, B: 0x4f, A: 0xff}
var Purple50 = &color.NRGBA{R: 0xf3, G: 0xe5, B: 0xf5, A: 0xff}
var Purple100 = &color.NRGBA{R: 0xe1, G: 0xbe, B: 0xe7, A: 0xff}
var Purple200 = &color.NRGBA{R: 0xce, G: 0x93, B: 0xd8, A: 0xff}
var Purple300 = &color.NRGBA{R: 0xba, G: 0x68, B: 0xc8, A: 0xff}
var Purple400 = &color.NRGBA{R: 0xab, G: 0x47, B: 0xbc, A: 0xff}
var Purple500 = &color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0xff}
var Purple600 = &color.NRGBA{R: 0x8e, G: 0x24, B: 0xaa, A: 0xff}
var Purple700 = &color.NRGBA{R: 0x7b, G: 0x1f, B: 0xa2, A: 0xff}
var Purple800 = &color.NRGBA{R: 0x6a, G: 0x1b, B: 0x9a, A: 0xff}
var Purple900 = &color.NRGBA{R: 0x4a, G: 0x14, B: 0x8c, A: 0xff}
var Indigo50 = &color.NRGBA{R: 0xe8, G: 0xea, B: 0xf6, A: 0xff}
var Indigo100 = &color.NRGBA{R: 0xc5, G: 0xca, B: 0xe9, A: 0xff}
var Indigo200 = &color.NRGBA{R: 0x9f, G: 0xa8, B: 0xda, A: 0xff}
var Indigo300 = &color.NRGBA{R: 0x79, G: 0x86, B: 0xcb, A: 0xff}
var Indigo400 = &color.NRGBA{R: 0x5c, G: 0x6b, B: 0xc0, A: 0xff}
var Indigo500 = &color.NRGBA{R: 0x3f, G: 0x51, B: 0xb5, A: 0xff}
var Indigo600 = &color.NRGBA{R: 0x39, G: 0x49, B: 0xab, A: 0xff}
var Indigo700 = &color.NRGBA{R: 0x30, G: 0x3f, B: 0x9f, A: 0xff}
var Indigo800 = &color.NRGBA{R: 0x28, G: 0x35, B: 0x93, A: 0xff}
var Indigo900 = &color.NRGBA{R: 0x1a, G: 0x23, B: 0x7e, A: 0xff}
var Blue50 = &color.NRGBA{R: 0xe3, G: 0xf2, B: 0xfd, A: 0xff}
var Blue100 = &color.NRGBA{R: 0xbb, G: 0xde, B: 0xfb, A: 0xff}
var Blue200 = &color.NRGBA{R: 0x90, G: 0xca, B: 0xf9, A: 0xff}
var Blue300 = &color.NRGBA{R: 0x64, G: 0xb5, B: 0xf6, A: 0xff}
var Blue400 = &color.NRGBA{R: 0x42, G: 0xa5, B: 0xf5, A: 0xff}
var Blue500 = &color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0xff}
var Blue600 = &color.NRGBA{R: 0x1e, G: 0x88, B: 0xe5, A: 0xff}
var Blue700 = &color.NRGBA{R: 0x19, G: 0x76, B: 0xd2, A: 0xff}
var Blue800 = &color.NRGBA{R: 0x15, G: 0x65, B: 0xc0, A: 0xff}
var Blue900 = &color.NRGBA{R: 0x0d, G: 0x47, B: 0xa1, A: 0xff}
var Cyan50 = &color.NRGBA{R: 0xe0, G: 0xf7, B: 0xfa, A: 0xff}
var Cyan100 = &color.NRGBA{R: 0xb2, G: 0xeb, B: 0xf2, A: 0xff}
var Cyan200 = &color.NRGBA{R: 0x80, G: 0xde, B: 0xea, A: 0xff}
var Cyan300 = &color.NRGBA{R: 0x4d, G: 0xd0, B: 0xe1, A: 0xff}
var Cyan400 = &color.NRGBA{R: 0x26, G: 0xc6, B: 0xda, A: 0xff}
var Cyan500 = &color.NRGBA{R: 0x00, G: 0xbc, B: 0xd4, A: 0xff}
var Cyan600 = &color.NRGBA{R: 0x00, G: 0xac, B: 0xc1, A: 0xff}
var Cyan700 = &color.NRGBA{R: 0x00, G: 0x97, B: 0xa7, A: 0xff}
var Cyan800 = &color.NRGBA{R: 0x00, G: 0x83, B: 0x8f, A: 0xff}
var Cyan900 = &color.NRGBA{R: 0x00, G: 0x60, B: 0x64, A: 0xff}
var Teal50 = &color.NRGBA{R: 0xe0, G: 0xf2, B: 0xf1, A: 0xff}
var Teal100 = &color.NRGBA{R: 0xb2, G: 0xdf, B: 0xdb, A: 0xff}
var Teal200 = &color.NRGBA{R: 0x80, G: 0xcb, B: 0xc4, A: 0xff}
var Teal300 = &color.NRGBA{R: 0x4d, G: 0xb6, B: 0xac, A: 0xff}
var Teal400 = &color.NRGBA{R: 0x26, G: 0xa6, B: 0x9a, A: 0xff}
var Teal500 = &color.NRGBA{R: 0x00, G: 0x96, B: 0x88, A: 0xff}
var Teal600 = &color.NRGBA{R: 0x00, G: 0x89, B: 0x7b, A: 0xff}
var Teal700 = &color.NRGBA{R: 0x00, G: 0x79, B: 0x6b, A: 0xff}
var Teal800 = &color.NRGBA{R: 0x00, G: 0x69, B: 0x5c, A: 0xff}
var Teal900 = &color.NRGBA{R: 0x00, G: 0x4d, B: 0x40, A: 0xff}
var Green50 = &color.NRGBA{R: 0xe8, G: 0xf5, B: 0xe9, A: 0xff}
var Green100 = &color.NRGBA{R: 0xc8, G: 0xe6, B: 0xc9, A: 0xff}
var Green200 = &color.NRGBA{R: 0xa5, G: 0xd6, B: 0xa7, A: 0xff}
var Green300 = &color.NRGBA{R: 0x81, G: 0xc7, B: 0x84, A: 0xff}
var Green400 = &color.NRGBA{R: 0x66, G: 0xbb, B: 0x6a, A: 0xff}
var Green500 = &color.NRGBA{R: 0x4c, G: 0xaf, B: 0x50, A: 0xff}
var Green600 = &color.NRGBA{R: 0x43, G: 0xa0, B: 0x47, A: 0xff}
var Green700 = &color.NRGBA{R: 0x38, G: 0x8e, B: 0x3c, A: 0xff}
var Green800 = &color.NRGBA{R: 0x2e, G: 0x7d, B: 0x32, A: 0xff}
var Green900 = &color.NRGBA{R: 0x1b, G: 0x5e, B: 0x20, A: 0xff}
var Yellow50 = &color.NRGBA{R: 0xff, G: 0xfd, B: 0xe7, A: 0xff}
var Yellow100 = &color.NRGBA{R: 0xff, G: 0xf9, B: 0xc4, A: 0xff}
var Yellow200 = &color.NRGBA{R: 0xff, G: 0xf5, B: 0x9d, A: 0xff}
var Yellow300 = &color.NRGBA{R: 0xff, G: 0xf1, B: 0x76, A: 0xff}
var Yellow400 = &color.NRGBA{R: 0xff, G: 0xee, B: 0x58, A: 0xff}
var Yellow500 = &color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0xff}
var Yellow600 = &color.NRGBA{R: 0xfd, G: 0xd8, B: 0x35, A: 0xff}
var Yellow700 = &color.NRGBA{R: 0xfb, G: 0xc0, B: 0x2d, A: 0xff}
var Yellow800 = &color.NRGBA{R: 0xf9, G: 0xa8, B: 0x25, A: 0xff}
var Yellow900 = &color.NRGBA{R: 0xf5, G: 0x7f, B: 0x17, A: 0xff}
var Orange50 = &color.NRGBA{R: 0xff, G: 0xf3, B: 0xe0, A: 0xff}
var Orange100 = &color.NRGBA{R: 0xff, G: 0xe0, B: 0xb2, A: 0xff}
var Orange200 = &color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff}
var Orange300 = &color.NRGBA{R: 0xff, G: 0xb7, B: 0x4d, A: 0xff}
var Orange400 = &color.NRGBA{R: 0xff, G: 0xa7, B: 0x26, A: 0xff}
var Orange500 = &color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff}
var Orange600 = &color.NRGBA{R: 0xfb, G: 0x8c, B: 0x00, A: 0xff}
var Orange700 = &color.NRGBA{R: 0xf5, G: 0x7c, B: 0x00, A: 0xff}
var Orange800 = &color.NRGBA{R: 0xef, G: 0x6c, B: 0x00, A: 0xff}
var Orange900 = &color.NRGBA{R: 0xe6, G: 0x51, B: 0x00, A: 0xff}
var Brown50 = &color.NRGBA{R: 0xef, G: 0xeb, B: 0xe9, A: 0xff}
var Brown100 = &color.NRGBA{R: 0xd7, G: 0xcc, B: 0xc8, A: 0xff}
var Brown200 = &color.NRGBA{R: 0xbc, G: 0xaa, B: 0xa4, A: 0xff}
var Brown300 = &color.NRGBA{R: 0xa1, G: 0x88, B: 0x7f, A: 0xff}
var Brown400 = &color.NRGBA{R: 0x8d, G: 0x6e, B: 0x63, A: 0xff}
var Brown500 = &color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0xff}
var Brown600 = &color.NRGBA{R: 0x6d, G: 0x4c, B: 0x41, A: 0xff}
var Brown700 = &color.NRGBA{R: 0x5d, G: 0x40, B: 0x37, A: 0xff}
var Brown800 = &color.NRGBA{R: 0x4e, G: 0x34, B: 0x2e, A: 0xff}
var Brown900 = &color.NRGBA{R: 0x3e, G: 0x27, B: 0x23, A: 0xff}
var Gray50 = &color.NRGBA{R: 0xfa, G: 0xfa, B: 0xfa, A: 0xff}
var Gray100 = &color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
var Gray200 = &color.NRGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff}
var Gray300 = &color.NRGBA{R: 0xe0, G: 0xe0, B: 0xe0, A: 0xff}
var Gray400 = &color.NRGBA{R: 0xbd, G: 0xbd, B: 0xbd, A: 0xff}
var Gray500 = &color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0xff}
var Gray600 = &color.NRGBA{R: 0x75, G: 0x75, B: 0x75, A: 0xff}
var Gray700 = &color.NRGBA{R: 0x61, G: 0x61, B: 0x61, A: 0xff}
var Gray800 = &color.NRGBA{R: 0x42, G: 0x42, B: 0x42, A: 0xff}
var Gray900 = &color.NRGBA{R: 0x42, G: 0x42, B: 0x42, A: 0xff}
var BlueGray50 = &color.NRGBA{R: 0xec, G: 0xef, B: 0xf1, A: 0xff}
var BlueGray100 = &color.NRGBA{R: 0xcf, G: 0xd8, B: 0xdc, A: 0xff}
var BlueGray200 = &color.NRGBA{R: 0xb0, G: 0xbe, B: 0xc5, A: 0xff}
var BlueGray300 = &color.NRGBA{R: 0x90, G: 0xa4, B: 0xae, A: 0xff}
var BlueGray400 = &color.NRGBA{R: 0x78, G: 0x90, B: 0x9c, A: 0xff}
var BlueGray500 = &color.NRGBA{R: 0x60, G: 0x7d, B: 0x8b, A: 0xff}
var BlueGray600 = &color.NRGBA{R: 0x54, G: 0x6e, B: 0x7a, A: 0xff}
var BlueGray700 = &color.NRGBA{R: 0x45, G: 0x5a, B: 0x64, A: 0xff}
var BlueGray800 = &color.NRGBA{R: 0x37, G: 0x47, B: 0x4f, A: 0xff}
var BlueGray900 = &color.NRGBA{R: 0x26, G: 0x32, B: 0x38, A: 0xff}
