package main

import (
	"net/http"
	"strings"
	//"time"
)

//функция для лога всех запросов
func LogReq(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//startLoadTime := time.Now()
		//LogPrint(CurTimeStrRFC3339() + " <- " + r.URL.Scheme + " " + r.URL.Path + "?" + r.URL.RawQuery)

		f(w, r)

		//LogPrint(CurTimeStrRFC3339() + " -> " + r.URL.Scheme + " " + r.URL.Path + "?" + r.URL.RawQuery + sprintf("  %v ", time.Now().Sub(startLoadTime)))
	}
}

func ShowError(title string, err error, w http.ResponseWriter, r *http.Request) {
	if err == nil {
		return
	}
	serr := "\n\n== ERROR: ======================================\n"
	serr += title + "\n"
	serr += ErrStr(err)
	serr += "\n\n== /ERROR ======================================\n"
	serr = strings.Replace(serr, "\n", "\n#", -1)
	LogPrint(serr)

	//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(serr))

}

//возвращает значение элемента d[args1][args2][args3]
//или возвращает defaultval если этот элемент не существует
func get_map_val(d map[string]interface{}, defaultval interface{}, args ...string) interface{} {
	var t interface{}
	t = d
	for _, v := range args {
		x := t.(map[string]interface{})
		a, ok := x[v]
		if !ok {
			return defaultval
		}
		t = a
	}
	return t
}
