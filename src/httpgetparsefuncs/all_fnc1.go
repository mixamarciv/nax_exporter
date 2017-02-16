package httpgetparsefuncs

import (
	"fmt"
	//"strconv"
	"time"

	"encoding/json"
	"log"

	"regexp"
	"runtime/debug"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var fnclist map[string](func(string) string)

func Get_parser_fnc(fncname string) func(string) string {
	f, b := fnclist[fncname]
	if !b {
		return func(string) string {
			return "#error parser func " + fncname + " not found!\nERROR_not_found_function_" + fncname + " 1"
		}
	}
	return f
}

var sprintf = fmt.Sprintf

func i64toa(d int64) string {
	return sprintf("%d", d)
}

func itoa(d int) string {
	return sprintf("%d", d)
}

//функции для работы с ошибками
func ErrStr(err error) string {
	s := fmt.Sprintf("%+v", err)
	a := string(debug.Stack())
	//убераем указатель на текущую строку
	i := strings.Index(a, "\n")
	a = a[i+1:]
	i = strings.Index(a, "\n")
	a = a[i+1:]

	s += "\n" + a
	return s
}

//функции вывода даты и времени
func CurTimeStr() string {
	t := time.Now()
	p := fmt.Sprintf("%s", strings.Replace(t.Format(time.RFC3339)[0:19], "T", " ", 1))
	return p
}

func CurTimeStrRFC3339() string {
	t := time.Now()
	p := t.Format(time.RFC3339)[0:19]
	return p
}

//возвращает 20160926-095323
func CurTimeStrShort() string {
	//2016-04-02T18:21:09+03:00
	t := time.Now()
	p := fmt.Sprintf("%s", t.Format(time.RFC3339)[0:19])
	p = p[0:19]
	p = strings.Replace(p, "-", "", -1)
	p = strings.Replace(p, ":", "", -1)
	p = strings.Replace(p, "T", "-", -1)
	return p
}

//функции для работы со строками
func RegexpCompile(re string) (*regexp.Regexp, error) {
	return regexp.Compile(re)
}

func StrRegexpMatch(re, s string) bool {
	r, err := regexp.Compile(re)
	if err != nil {
		//printerr("RegexpMatch Compile error", err)
		log.Panicln("RegexpMatch Compile("+re+") error", err)
	}
	return r.MatchString(s)
}
func StrRegexpReplace(text string, regx_from string, to string) string {
	reg, err := regexp.Compile(regx_from)
	if err != nil {
		log.Panicln("StrRegexpReplace Compile("+regx_from+") error", err)
	}
	text = reg.ReplaceAllString(text, to)
	return text
}

//преобразует из json строки в map[string]interface{}
func FromJsonStr(data []byte) map[string]interface{} {
	var d map[string]interface{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return map[string]interface{}{"_json_parse_error": "1", "error": ErrStr(err), "data": string(data)}
	}
	return d
}
func ToJsonStr(v interface{}) string {
	j, err := json.Marshal(v)
	if err != nil {
		return ErrStr(err)
	}
	return string(j)
}

func var_dump(v interface{}, depth int, indent string) string {
	cs := &spew.ConfigState{
		Indent:                  indent,
		MaxDepth:                depth,
		SortKeys:                true,
		DisableMethods:          true,
		DisableCapacities:       true,
		DisablePointerAddresses: true,
		DisablePointerMethods:   true,
		SpewKeys:                true,
	}
	//return cs.Sprintf("%#v", v)
	return cs.Sdump(v)
}
