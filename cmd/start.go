package cmd

import (
	"fmt"
	"os"
	"strings"
	"transaction-matching-engine/engine"
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
			os.Exit(0)
		},
	}

	startGrpcCmd = &cobra.Command{
		Use:    "grpc",
		Short:  "启动grpc服务",
		PreRun: preServerRun,
		Run: func(cmd *cobra.Command, args []string) {
			var pairs []string
			if len(args) > 0 {
				pairs = strings.Split(strings.ToUpper(args[0]), ",")
			} else {
				pairs = engine.ReadPairs()
				if len(pairs) == 0 {
					panic("转储的交易对文件不存在或无交易对")
				}
			}
			engine.Load(pairs)
			grpc.Run(pairs)
		},
	}
)

// 服务启动前参数检查，数据载入
func preServerRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("未在启动参数指定交易对,例： [start x服务 BTC-USDT] \r\n交易对将从配置文件加载...")
	}
}
