package main

import (
	"fmt"
	"os"
	"runtime"
	s "strings"
)

var log_file string

func InitLog() {
	log_path, _ := AppPath()
	if runtime.GOOS == "windows" {
		log_path = s.Replace(log_path, "\\", "/", -1)
	} else { //if runtime.GOOS == "linux" {
		log_path = "/var/log"
	}

	timestr := CurTimeStrShort()
	log_path = log_path + "/log_nax_exporter/" + timestr[0:8]

	MkdirAll(log_path)
	log_file = log_path + "/" + CurTimeStrShort() + ".log"

	WriteLogln("start log")
	WriteLogln("log file: " + log_file)
}

func WriteLog(data string) {
	FileAppendStr(log_file, data)
}

func WriteLogln(data string) {
	FileAppendStr(log_file, CurTimeStr()+" "+s.TrimRight(data, "\n\r\t ")+"\n")
}

func WriteLogErr(info string, err error) {
	FileAppendStr(log_file, CurTimeStr()+" "+info+"\n"+ErrStr(err))
}

func WriteLogErrAndExit(info string, err error) {
	if err == nil {
		return
	}
	FileAppendStr(log_file, CurTimeStr()+" "+info+"\n"+ErrStr(err))
	panic(err)
	os.Exit(1)
}

func LogPrint(data string) {
	fmt.Println(data)
	WriteLogln(data)
}

func LogPrintErrAndExit(info string, err error) {
	if err == nil {
		return
	}
	fmt.Println(info)
	fmt.Printf("%+v", err)
	WriteLogErrAndExit(info, err)
}

func LogPrintErr(info string, err error) {
	if err == nil {
		return
	}
	fmt.Println(info)
	fmt.Printf("%+v", err)
	WriteLogErr(info, err)
}

func LogPrintAndExit(info string) {
	fmt.Println(info)
	WriteLogln(info)
	os.Exit(1)
}
