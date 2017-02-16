package main

import (
	"net/http"
	"net/url"
	//"os/exec"
	//"regexp"
	//"runtime"
	f "httpgetparsefuncs"
	"strings"
)

var httpget_workers_cnt = 4

func init() {
	httpget_workers_cnt = int(gcfg_app["httpget_workers_cnt"].(float64))
}

func exp_httpget(w http.ResponseWriter, r *http.Request) {
	GET, _ := url.ParseQuery(r.URL.RawQuery)

	urls, b := GET["url"]
	if !b {
		w.Write([]byte("GET['url'] not set\n"))
		return
	}

	vtype, b := GET["type"]
	if !b {
		w.Write([]byte("GET['type'] not set\n"))
		return
	}
	vtype1 := vtype[0]

	cnt_j := len(urls)           //количество заданий
	cnt_w := httpget_workers_cnt //максимальное количество горутин
	if cnt_j < cnt_w {           //ограничеваем количество горутин
		cnt_w = cnt_j //должно быть не больше чем всего заданий
	}

	jobs := make(chan string, cnt_w)
	results := make(chan string, cnt_w)
	for w := 1; w <= cnt_w; w++ {
		go exp_httpget_worker(w, jobs, results) //стартуем горутины
	}

	for j := 0; j < cnt_j; j++ {
		jobs <- vtype1 + "|" + urls[j]
	}
	close(jobs)

	data := ""
	for r := 0; r < cnt_j; r++ {
		data += <-results
		data += "\n"
	}

	w.Write([]byte(data))
}

func exp_httpget_worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		s := sprintf("#worker%d %s\n", id, j)
		s += exp_httpget_url(j)
		results <- s
	}
}

func exp_httpget_url(job string) string {
	a := strings.Split(job, "|")
	vtype1 := a[0]
	url := a[1]

	//data := sprintf("#%s type: %s url: %s\n", f.CurTimeStrShort(), vtype1, url)

	body, err := SendHttpRequest(url, map[string]string{})

	fnc := f.Get_parser_fnc(vtype1)
	data := ""
	if err != nil {
		str := "#" + ErrStr(err)
		str = strings.Replace(str, "\n", "\n#", -1)
		data += str
	} else {
		data = fnc(string(body))
	}

	return data
}
