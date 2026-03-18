package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "base64:decode",
		Short: "Base64 解码",
		Run:   runBase64Decode,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})

	appendCommand(&cobra.Command{
		Use:   "base64:encode",
		Short: "Base64 编码",
		Run:   runBase64Encode,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runBase64Decode(cmd *cobra.Command, args []string) {
	var input []byte
	if len(args) > 0 {
		input = []byte(args[0])
	} else {
		var err error
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			return
		}
	}
	res, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}
}

func runBase64Encode(cmd *cobra.Command, args []string) {
	var input []byte
	if len(args) > 0 {
		input = []byte(args[0])
	} else {
		var err error
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			return
		}
	}

	res := base64.StdEncoding.EncodeToString(input)
	fmt.Println(string(res))
}
