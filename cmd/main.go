package main

import (
	"os"
	"strconv"
	. "github.com/pastet89/alien-invasion/utils"
	. "github.com/pastet89/alien-invasion/simulator"
)

func main() {
	simulator := ConstructSimulator(getInputs())
	simulator.Run()
}

func getInputs() (int, string) {
	aliensNumStr, MapFilePath := os.Args[1], os.Args[2]
	aliensNum, err := strconv.Atoi(aliensNumStr)
	ProcessError(err)
	return aliensNum, MapFilePath
}