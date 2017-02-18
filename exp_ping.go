package main

import (
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var re_not_d *regexp.Regexp
var re_not_d2 *regexp.Regexp
var ping_workers_cnt = 4

func init() {
	re_not_d, _ = regexp.Compile("[^\\d]+")
	re_not_d2, _ = regexp.Compile("[^\\d\\<\\.]+")
	ping_workers_cnt = int(gcfg_app["ping_workers_cnt"].(float64))
}

func exp_ping(w http.ResponseWriter, r *http.Request) {
	GET, _ := url.ParseQuery(r.URL.RawQuery)

	hosts, b := GET["host"]
	if !b {
		w.Write([]byte("GET['host'] not set\n"))
		return
	}

	cnt := len(hosts)
	cnt_w := ping_workers_cnt
	if cnt < cnt_w {
		cnt_w = cnt
	}

	jobs := make(chan string, cnt)
	results := make(chan string, cnt)
	for w := 1; w <= cnt_w; w++ {
		go exp_ping_worker(w, jobs, results)
	}

	for j := 0; j < cnt; j++ {
		jobs <- hosts[j]
	}
	close(jobs)

	data := ""
	for r := 0; r < cnt; r++ {
		data += <-results
		data += "\n"
	}

	w.Write([]byte(data))
}

func exp_ping_worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		s := sprintf("#worker %d\n", id)
		s += exp_ping_host(j)
		results <- s
	}
}

func exp_ping_host(host string) string {
	var cmd []string
	if runtime.GOOS == "windows" {
		//cmd = []string{"chcp", "65001", "&&", "ping", host, "-n", "1"}
		//cmd = []string{"cmd", "/c", "\"chcp 65001&ping " + host + " -n 1\""}
		cmd = []string{"ping", host, "-n", "1", "-w", "2000"}
	} else { //if runtime.GOOS == "linux" {
		cmd = []string{"ping", host, "-c", "1", "-W", "2"}
	}
	//w.Write([]byte(sprintf("#cmd: %#v\n\n", cmd)))

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	s := string(out)

	ms := "90000"
	ttl := "0"
	success := "0"

	if runtime.GOOS == "windows" {
		//Ответ от 192.168.0.30: число байт=32 время<1мс TTL=128
		i := strings.Index(s, "TTL=")
		if i > 0 {
			success = "1"

			//ttl = StrRegexpReplace(s[i+4:i+7], "[^\\d]+", "")
			ttl = re_not_d.ReplaceAllString(s[i+4:i+7], "")

			//ms = StrRegexpReplace(ms, "[^\\d\\<\\.]+", "")
			ms = re_not_d2.ReplaceAllString(s[i-10:i], "")
			if ms[0] == '<' {
				ms = "0.5"
			}
		}
	} else {
		//64 bytes from 192.168.0.30: icmp_seq=1 ttl=128 time=0.284 ms
		i := strings.Index(s, "ttl=")
		if i > 0 {
			success = "1"

			ttl = strings.TrimRight(s[i+4:i+7], " ")

			ms = s[i+7:]
			i = strings.Index(ms, "=")
			j := strings.Index(ms[i:], " ")
			ms = ms[i+1 : i+j]
		}
	}

	data := "ping_ms{ip=\"" + host + "\"} " + ms + "\n"
	data += "ping_ttl{ip=\"" + host + "\"} " + ttl + "\n"
	data += "ping_success{ip=\"" + host + "\"} " + success + "\n"

	//LogPrint("out: " + string(out))
	if err != nil {
		data += "ping_error{ip=\"" + host + "\"} 1\n"
		data += sprintf("#cmd: %#v\n#", cmd)
		str := ErrStr(err)
		str = strings.Replace(str, "\n", "\n#", -1)
		data += str
	}
	//w.Write([]byte(data))
	return data
}
