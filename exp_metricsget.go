package main

import (
	"net/http"
	"net/url"
	//"os/exec"
	//"regexp"
	//"runtime"
	a "app_fnc"
	f "metricsgetfuncs"
	"strings"
)

var metricsget_workers_cnt = 4
var metricsget_params_sep = "-||@"

func init() {
	metricsget_workers_cnt = int(gcfg_app["metricsget_workers_cnt"].(float64))
}

func exp_metricsget(w http.ResponseWriter, r *http.Request) {
	GET, _ := url.ParseQuery(r.URL.RawQuery)

	cnt_w := httpget_workers_cnt //максимальное количество горутин
	threads, b := GET["threads"]
	if b {
		need_cnt_w, err := a.Atoi(threads[0])
		if err != nil {
			ret := a.ErrStr2Comment("bad threads number: "+threads[0], err)
			w.Write([]byte(ret))
			return
		}
		if need_cnt_w >= 1 {
			cnt_w = need_cnt_w
		}
	}

	//список функций
	fncs, b := GET["fnc"]
	if !b {
		w.Write([]byte("GET['fnc'] not set\n"))
		return
	}

	GET["__hash_rawquery"] = []string{a.StrCrc32([]byte(r.URL.RawQuery))}
	GET["__cnt_w"] = []string{a.Itoa(cnt_w)}

	//запускаем воркеров
	res := a.Exec_go_workers(exp_metricsget_fnc, fncs, GET, cnt_w)

	data := strings.Join(res, "\n\n")

	w.Write([]byte(data))
}

func exp_metricsget_fnc(job string, get_params interface{}) string {
	fnc_name := job
	fnc := f.Get_parser_fnc(fnc_name)
	data := fnc(fnc_name, get_params)

	return data
}
