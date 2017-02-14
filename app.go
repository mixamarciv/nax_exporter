package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	InitLog()
	InitApp()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", LogReq(exp_ping))

	srv := &http.Server{
		Handler:      r,
		Addr:         gcfg_app["server_addr"].(string),
		WriteTimeout: 400 * time.Second,
		ReadTimeout:  400 * time.Second,
	}

	LogPrint("start listening addr: " + gcfg_app["server_addr"].(string))
	err := srv.ListenAndServe()
	LogPrintErrAndExit("ERROR start listening addr: "+gcfg_app["server_addr"].(string)+" ", err)
}
