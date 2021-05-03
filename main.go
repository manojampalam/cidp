package main

import (
	"os"

	getopt "github.com/pborman/getopt/v2"
)

var (
	argsConfigFile = ""
	testMode       bool
	help           bool
)

func init() {
	getopt.FlagLong(&argsConfigFile, "config", 'c', "path to k8s config file")
	getopt.Flag(&testMode, 't', "run in test mode")
	getopt.Flag(&help, '?', "display help")
}

func main() {
	getopt.SetParameters("")
	getopt.Parse()

	if help == true {
		getopt.PrintUsage(os.Stdout)
		return
	}

	run(argsConfigFile)
}
