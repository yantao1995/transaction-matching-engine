package main

import (
	"transaction-matching-engine/cmd"
	"transaction-matching-engine/common"
)

func main() {
	cmd.Execute()
	common.ServerStatus.Wait()
}
