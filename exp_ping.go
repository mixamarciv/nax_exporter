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

func init() {
	re_not_d, _ = regexp.Compile("[^\\d]+")
	re_not_d2, _ = regexp.Compile("[^\\d\\<\\.]+")
}

func exp_ping(w http.ResponseWriter, r *http.Request) {
	GET, _ := url.ParseQuery(r.URL.RawQuery)

	t, b := GET["host"]
	if !b {
		w.Write([]byte("GET['host']\n"))
	}
	host := t[0]

	var cmd []string
	if runtime.GOOS == "windows" {
		//cmd = []string{"chcp", "65001", "&&", "ping", host, "-n", "1"}
		//cmd = []string{"cmd", "/c", "\"chcp 65001&ping " + host + " -n 1\""}
		cmd = []string{"ping", host, "-n", "1"}
	} else { //if runtime.GOOS == "linux" {
		cmd = []string{"ping", host, "-c", "1"}
	}
	//w.Write([]byte(sprintf("#cmd: %#v\n\n", cmd)))

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	s := string(out)

	ms := "-1"
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

	data := "ping_" + host + "_ms " + ms + "\n"
	data += "ping_" + host + "_ttl " + ttl + "\n"
	data += "ping_" + host + "_success " + success + "\n\n"

	//LogPrint("out: " + string(out))
	if err != nil {
		data += "ping_" + host + "_error 1\n\n"
		data += sprintf("#cmd: %#v\n\n#", cmd)
		str := ErrStr(err)
		str = strings.Replace(str, "\n", "\n#", -1)
		data += str
	}
	w.Write([]byte(data))
}
