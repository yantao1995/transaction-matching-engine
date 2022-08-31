package main

import (
	"transaction-matching-engine/cmd"
	"transaction-matching-engine/common"
	"transaction-matching-engine/engine"
)

func main() {
	cmd.Execute()
	common.ServerStatus.Wait()
	engine.Dump()
}
