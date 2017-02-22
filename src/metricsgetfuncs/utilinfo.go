package metricsgetfuncs

import (
	a "app_fnc"
	"bytes"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

/***********
  #запуск утилит и парсинг результатов
  - job_name: 'utilinfo'
    scrape_interval: 60s
    metrics_path: '/metricsget'
    scheme: 'http'
    params:
        threads: [1,2]                      # сколько функций будет выполняться одновременно, сколько утилит будет запущено одновременно default[1,1]
        fnc: ['utilinfo']                   # список функций
        util: ['top','who']                 # список утилит
        other_params: ['asdf1234','789']    # другие параметры - доступны в каждой функции
    static_configs:
      - targets: ['192.168.0.31:1002']
************/
type utilfnc func(string, map[string]interface{}) string

//список функций - для запуска и парсинга утилит
var util_list map[string]utilfnc

func init() {
	util_list = make(map[string]utilfnc)
	util_list["top"] = utilinfo__util_top
	util_list["who"] = utilinfo__util_who
}

func utilinfo(fnc_name string, params interface{}) string {

	GET := params.(url.Values)
	d, errstr := utilinfo__get_opt(GET)
	if errstr != "" {
		return errstr
	}

	cnt_w := d["cnt_w"].(int)
	utils := d["utils"].([]string)

	res := a.Exec_go_workers(RunAndParseUtil, utils, d, cnt_w)

	ret := "#" + d["info"].(string) + "\n"
	ret += strings.Join(res, "\n\n")

	return ret
}

//возвращает список урл и регулярных выражений для проверки
func utilinfo__get_opt(GET url.Values) (map[string]interface{}, string) {
	hash := GET["__hash_rawquery"][0]
	d, b := cache_opt[hash]
	if b {
		d["__last_update"] = time.Now()
		return d, ""
	}

	d = make(map[string]interface{})

	utils, b := GET["util"]
	fnclist := make(map[string]utilfnc)
	if !b {
		return nil, "GET['util'] not set\n"
	}
	for _, fname := range utils {
		f, b := util_list[fname]
		if !b {
			f = utilinfo__util_not_found
		}
		fnclist[fname] = f
	}

	//определяем количество воркеров которые будут проверять список урл
	cnt_w := 8
	tcnt_w, b := GET["threads"]
	if b && len(tcnt_w) >= 2 {
		i, err := a.Atoi(tcnt_w[1])
		if err == nil && i > 1 {
			cnt_w = i
		}
	}

	d["fnclist"] = fnclist
	d["utils"] = utils
	d["cnt_w"] = cnt_w

	d["info"] = "threads:" + a.Sprintf("%v,%d", GET["__cnt_w"][0], cnt_w)

	return d, ""
}

func RunAndParseUtil(fname string, p interface{}) string {
	d := p.(map[string]interface{})
	fnclist := d["fnclist"].(map[string]utilfnc)
	f := fnclist[fname]
	s := f(fname, d)
	return s
}

func utilinfo__util_not_found(fname string, d map[string]interface{}) string {
	return "ERROR_nax_exporter{utilinfo_not_found_util_function=\"" + fname + "\"} 1"
}

func utilinfo__util_top(fname string, d map[string]interface{}) string {
	var cmd []string
	if runtime.GOOS == "windows" {
		return "ERROR_nax_exporter{utilinfo_top_error=\"top_for_windows_not_found\"} 1"
	} else {
		cmd = []string{"top", "-n", "1", "-b"}
	}

	var out, outerr bytes.Buffer
	{
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Stdout = &out
		c.Stderr = &outerr
		err := c.Run()
		if err != nil {
			return a.ErrStr2Comment("exec.Command: "+strings.Join(cmd, " ")+"\nerror out:\n"+outerr.String()+"\nout:\n"+out.String(), err)
		}
	}
	s := out.String()

	ret_s := ""
	{ //user
		//top - 19:28:23 up 19:27,  1 user,  load average: 0,00, 0,00, 0,00
		i := strings.Index(s, " user")
		if i < 0 {
			return "#not found string \" user\"\nERROR_nax_exporter{utilinfo_top_error=\"err5\"} 1"
		}

		j := strings.Index(s, ", ")
		if j < 0 {
			return "#not found string \", \"\nERROR_nax_exporter{utilinfo_top_error=\"err6\"} 1"
		}
		if j > i {
			return "#bad string\nERROR_nax_exporter{utilinfo_top_error=\"err7\"} 1"
		}

		t := s[j+2 : i]
		ret_s += "utilinfo_top_users " + t
	}

	{ //load average:
		//top - 19:28:23 up 19:27,  1 user,  load average: 0,00, 0,05, 0,25
		i := strings.Index(s, "load average:")
		if i < 0 {
			return "#not found string \"load average:\"\nERROR_nax_exporter{utilinfo_top_error=\"err2\"} 1" +
				"\n#out string:\n#" + strings.Replace(s, "\n", "\n#", -1) + "\n"
		}

		t := s[i+14 : i+50]
		i = strings.Index(t, "\n")
		if i < 0 {
			return "#bad string \"load average:\"\nERROR_nax_exporter{utilinfo_top_error=\"err3\"} 1"
		}

		t = t[:i]

		t = strings.Replace(t, ", ", " ", 2)
		t = strings.Replace(t, ",", ".", 3)

		a := strings.Split(t, " ")

		if len(a) < 3 {
			return "#bad string \"load average:\" result: \"" + t + "\"\nERROR_nax_exporter{utilinfo_top_error=\"err4\"} 1"
		}
		ret_s += "\nutilinfo_top_load_average_1 " + a[0] +
			"\nutilinfo_top_load_average_2 " + a[1] +
			"\nutilinfo_top_load_average_3 " + a[2]
	}

	{ //Tasks:
		//Tasks:  84 total,   1 running,  83 sleeping,   0 stopped,   0 zombie
		i := strings.Index(s, "Tasks:")
		if i < 0 {
			return "#not found string \"Tasks:\"\nERROR_nax_exporter{utilinfo_top_error=\"err8\"} 1"
		}

		j := strings.Index(s[i:i+200], "\n")
		if j < 0 {
			return "#not found string \"\n\"\nERROR_nax_exporter{utilinfo_top_error=\"err9\"} 1"
		}
		t := s[i+6 : i+j]
		a := strings.Split(t, ",")

		if len(a) < 4 {
			return "#bad string \"Tasks:\" result: \"" + t + "\"\nERROR_nax_exporter{utilinfo_top_error=\"err10\"} 1"
		}
		ret_s += "\nutilinfo_top_tasks_total " + a[0][:4] +
			"\nutilinfo_top_tasks_running " + a[1][:4] +
			"\nutilinfo_top_tasks_sleeping " + a[2][:4] +
			"\nutilinfo_top_tasks_zombie " + a[3][:4]
	}

	return ret_s
}

//who
func utilinfo__util_who(fname string, d map[string]interface{}) string {
	var cmd []string
	if runtime.GOOS == "windows" {
		return "ERROR_nax_exporter{utilinfo_who_error=\"who_for_windows_not_found\"} 1"
	} else {
		//cmd = []string{"who", "-a"}
		cmd = []string{"w", "-h"}
	}

	var out, outerr bytes.Buffer
	{
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Stdout = &out
		c.Stderr = &outerr
		err := c.Run()
		if err != nil {
			return a.ErrStr2Comment("exec.Command: "+strings.Join(cmd, " ")+"\nerror out:\n"+outerr.String()+"\nout:\n"+out.String(), err)
		}
	}
	s := out.String()
	s = strings.Replace(s, "    ", " ", -1)
	s = strings.Replace(s, "   ", " ", -1)
	s = strings.Replace(s, "  ", " ", -1)
	sr := strings.Split(s, "\n")

	ret_s := ""
	cnt_users := 0
	for _, s := range sr {
		//s == luser pts/0 192.168.0.30     13:40    4.00s  0.65s  0.05s sshd: luser [priv]
		a := strings.Split(s, " ")
		if len(a) < 3 {
			continue
		}
		cnt_users++
		ret_s += "utilinfo_who{user=\"" + a[0] + "\",tty=\"" + a[1] + "\",ip=\"" + a[2] + "\"} 1\n"
	}
	ret_s += "utilinfo_who_users " + a.Sprintf("%d", cnt_users)

	return ret_s
}
