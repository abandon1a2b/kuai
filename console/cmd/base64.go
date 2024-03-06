package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "base64:decode",
		Short: "文件读取",
		Run:   runBase64Decode,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})

	appendCommand(&cobra.Command{
		Use:   "base64:encode",
		Short: "文件读取",
		Run:   runBase64Encode,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runBase64Decode(cmd *cobra.Command, args []string) {
	buf := strings.Builder{}
	if len(args) > 0 {
		buf.WriteString(args[0])
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			result, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			buf.Write(result)
		}
	}
	res, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}
}

func runBase64Encode(cmd *cobra.Command, args []string) {

	buf := bytes.Buffer{}
	if len(args) > 0 {
		buf.WriteString(args[0])
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			result, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			buf.Write(result)
		}
	}

	res := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println(string(res))
}
