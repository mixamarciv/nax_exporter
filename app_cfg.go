package main

import (
	"os"
	"strings"
)

var apppath string
var gcfg map[string]string
var gcfg_app map[string]interface{}

func InitApp() {
	apppath, _ = AppPath()
	apppath = strings.Replace(apppath, "\\", "/", -1)

	if len(os.Args) < 2 {
		LogPrintAndExit("config file not set\nuse " + os.Args[0] + " config_file.json")
	}

	LogPrint("you run: " + apppath + "/" + os.Args[0] + " " + os.Args[1])

	file := os.Args[1]
	var err error
	gcfg_app, err = JsonFromFile(file)
	LogPrintErrAndExit("InitApp Error1 JsonFromFile: "+file+"\n\n", err)
}
