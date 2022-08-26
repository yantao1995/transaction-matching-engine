package cmd

import (
	"errors"
	"fmt"
	"strings"
	"transaction-matching-engine/grpc"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(start)
	start.AddCommand(startGrpcCmd)
}

var (
	start = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	startGrpcCmd = &cobra.Command{
		Use:     "grpc",
		Short:   "启动grpc服务",
		PreRunE: preServerRun,
		Run: func(cmd *cobra.Command, args []string) {
			pairs := strings.Split(args[0], ",")
			grpc.Run(pairs)
		},
		PostRun: postServerRun,
	}
)

//服务启动前参数检查，数据载入
func preServerRun(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	if len(args) < 1 {
		return errors.New("缺少启动参数，需要传入交易对,使用英文逗号(,)分隔,例如\r\nBTC-USDT,ETH-USDT")
	}
	return nil
}

//数据 dump
func postServerRun(cmd *cobra.Command, args []string) {
	fmt.Println("post run")
}
