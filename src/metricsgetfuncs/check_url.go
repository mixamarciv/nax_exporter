package metricsgetfuncs

import (
	a "app_fnc"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

/***********
  # проверка доступности сайтов
  - job_name: 'sites1'
    scrape_interval: 60s
    metrics_path: '/metricsget'
    scheme: 'http'
    params:
        threads: [1,2]                          # сколько функций будет выполняться одновременно, сколько урл будет проверяться одновременно default[1,1]
        fnc: ['check_url']                  # список функций
        url: ['http://vrashke.ru','http://vrashke.net','http//anykeyadmin.info']  # список параметров-урл, - все параметры будут передаваться во все функции

	    # регулярные выражения применяются только в функции 'check_url' и проверяются все урл
	    re_success: ['mixamarciv','<nav']      # body урл должна соответствовать всем регулярным выражениям для успешной проверки
	    re_fail: []                            # body урл должна соответствовать хотябы одному регулярному выражению для неудачи
	    re_status_success: []                  # status урл должна соответствовать всем регулярным выражениям для успешной проверки
	    re_status_fail: []                     # status урл должна соответствовать хотябы одному регулярному выражению для неудачи

    static_configs:
      - targets: ['192.168.0.31:1000']
************/
func check_url(fnc_name string, params interface{}) string {
	//re := regexp.MustCompile("(?is)<body[^<]*<")
	//t0 := time.Now()

	GET := params.(url.Values)
	d, errstr := check_url__get_opt(GET)
	if errstr != "" {
		return errstr
	}

	cnt_w := d["cnt_w"].(int)
	urls := d["urls"].([]string)

	/****
		res := ""
		for i := 0; i < len(urls); i++ {
			res += check_url__item(urls[i], d)
		}
	***/

	res := a.Exec_go_workers(check_url__item, urls, d, cnt_w)

	ret := "#" + d["info"].(string) + "\n"
	ret += strings.Join(res, "\n")
	//ret += "\ncheck_url_total_ms " + a.Sprintf("%d", int64(time.Now().Sub(t0).Nanoseconds()/1000000)) + "\n"
	return ret
}

//возвращает список урл и регулярных выражений для проверки
func check_url__get_opt(GET url.Values) (map[string]interface{}, string) {
	hash := GET["__hash_rawquery"][0]
	d, b := cache_opt[hash]
	if b {
		d["__last_update"] = time.Now()
		return d, ""
	}

	d = make(map[string]interface{})

	urls, b := GET["url"]
	if !b {
		return nil, "GET['url'] not set\n"
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

	//определяем регулярные выражения
	re := make(map[string][]*regexp.Regexp)
	for _, re_type := range []string{"re_success", "re_fail", "re_status_success", "re_status_fail"} {
		rt, errstr := check_url__get_regexp__s2r(GET[re_type])
		if errstr != "" {
			return nil, errstr
		}
		re[re_type] = rt
	}

	d["urls"] = urls
	d["re"] = re
	d["cnt_w"] = cnt_w

	d["info"] = "threads:" + a.Sprintf("%v,%d", GET["__cnt_w"][0], cnt_w)

	return d, ""
}

//возвращает список регулярных выражений для проверки результата урл
func check_url__get_regexp__s2r(s []string) ([]*regexp.Regexp, string) {
	if len(s) == 0 {
		return nil, ""
	}
	r := make([]*regexp.Regexp, len(s))
	for i := 0; i < len(s); i++ {
		re, err := regexp.Compile(s[i])
		if err != nil {
			return nil, a.ErrStr2Comment("bad regexp: "+s[i], err)
		}
		r[i] = re
	}
	return r, ""
}

//проверка урл
/***********
в d["re"] ==
	re_success: [] # body урл должна соответствовать всем регулярным выражениям для успешной проверки
	re_fail: [] # body урл должна соответствовать хотябы одному регулярному выраженияю для неудачи
	re_status_success: [] # status урл должна соответствовать всем регулярным выражениям для успешной проверки
	re_status_fail: [] # status урл должна соответствовать хотябы одному регулярному выраженияю для неудачи
************/
func check_url__item(url string, p interface{}) string {
	d := p.(map[string]interface{})

	t0 := time.Now()
	res, err := a.SendHttpRequest(url)

	res_time := int64(time.Now().Sub(t0).Nanoseconds() / 1000000)
	ms := a.Sprintf("%d", res_time)

	ret := ""

	if err != nil {
		ret += "check_url_fail{url=\"" + url + "\"} 1\n"
		ret += "check_url_fail_ms{url=\"" + url + "\"} " + ms + "\n"
		ret += a.ErrStr2Comment("ERROR SendHttpRequest(): "+url, err)
		return ret
	}

	resp := res["response"].(*http.Response)
	stat := resp.Status // string // e.g. "200 OK"
	body := res["body"].(string)

	//теперь проверяем регулярными выражениями результаты
	re := d["re"].(map[string][]*regexp.Regexp)

	success := check_url__item_re(stat, body, re)

	contentlength := resp.ContentLength
	if contentlength < 0 {
		contentlength = int64(len(body))
	}

	if success == 0 {
		ret += "check_url_fail{url=\"" + url + "\"} 1\n"
		ret += "check_url_fail_ms{url=\"" + url + "\"} " + ms + "\n"
		ret += "check_url_fail_contentlength{url=\"" + url + "\"} " + a.Sprintf("%d", contentlength) + "\n"
		ret += body
	} else {
		ret += "check_url_success{url=\"" + url + "\"} 1\n"
		ret += "check_url_success_ms{url=\"" + url + "\"} " + ms + "\n"
		ret += "check_url_success_contentlength{url=\"" + url + "\"} " + a.Sprintf("%d", contentlength) + "\n"
	}

	return ret
}

//проверка статуса и урл регулярными выражениями
// 1 - успешно 0 - нет
func check_url__item_re(stat, body string, re map[string][]*regexp.Regexp) int {
	for _, r := range re["re_status_success"] {
		if r.MatchString(stat) != true {
			return 0
		}
	}
	for _, r := range re["re_status_fail"] {
		if r.MatchString(stat) == true {
			return 0
		}
	}

	for _, r := range re["re_success"] {
		if r.MatchString(body) != true {
			return 0
		}
	}
	for _, r := range re["re_fail"] {
		if r.MatchString(body) == true {
			return 0
		}
	}
	return 1
}
