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

	gcfg_app = make(map[string]interface{})
	//file := apppath + "/" + os.Args[1]
	file := os.Args[1]
	data, err := FileRead(file)
	LogPrintErrAndExit("InitAppCfg Error1: can't read file: "+file+"\n\n", err)
	gcfg_app, err = FromJson(data)
	LogPrintErrAndExit("InitAppCfg Error2: Unmarshal json error: "+file+"\n\n", err)
}
