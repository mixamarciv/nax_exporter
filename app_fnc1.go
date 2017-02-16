package main

import (
	"encoding/base64"
	"fmt"
	"os"
	//"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var sprintf = fmt.Sprintf

func i64toa(d int64) string {
	return sprintf("%d", d)
}

func itoa(d int) string {
	return sprintf("%d", d)
}

func floatToStr(f interface{}) string {
	return strconv.FormatFloat(f.(float64), 'f', 0, 64)
}

func fmtError(s string, err error) string {
	return s + fmt.Sprintf("\n\n%#v", err)
}

func base64Decode(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Sprint("error:", err)
	}
	return string(decoded)
}

//текущий путь к приложению
func AppPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, err
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

//преобразует из json строки в map[string]interface{}
func FromJson(data []byte) (map[string]interface{}, error) {
	var d map[string]interface{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return map[string]interface{}{"error": ErrStr(err), "data": string(data)}, err
	}
	return d, nil
}

//функции для работы с файлами
func FileRead(file string) ([]byte, error) {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func FileReadStr(file string) (string, error) {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func FileWrite(file string, data []byte) error {
	err := ioutil.WriteFile(file, data, 0644)
	return err
}

func FileWriteStr(file string, data string) error {
	err := ioutil.WriteFile(file, []byte(data), 0644)
	return err
}

func FileAppendStr(filename string, data string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Panicln("FileAppendStr OpenFile error", err)
		//return err
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		log.Panicln("FileAppendStr WriteString error", err)
		//return err
	}
	f.Sync()
	return nil
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

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func MkdirAll(path string) error {
	return os.MkdirAll(path, 0777)
}

//отправляем запрос на список адресов по дому
func SendHttpRequest(urlStr string, m map[string]string) ([]byte, error) {
	client := &http.Client{}

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return /*[]byte{}*/ nil, err
	}

	//r.Header.Add("authority", "vk.com")
	//r.Header.Add("method", "GET")
	//r.Header.Add("path", "/")
	//r.Header.Add("scheme", "https")
	//r.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	//r.Header.Add("accept-encoding", "gzip, deflate, sdch, br")
	//r.Header.Add("accept-language", "ru-RU,ru;q=0.8,en-US;q=0.6,en;q=0.4")
	//r.Header.Add("upgrade-insecure-requests", "1")
	//r.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")

	resp, err := client.Do(r)
	if err != nil {
		return /*[]byte{}*/ nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data, err := UnzipData(body)
	if err == nil {
		body = data
	} else {
		//hz
	}

	//dump_resp := var_dump(resp, 4, "\t")

	//body = []byte(dump_resp + "\n\n\n" + string(body))

	return body, nil
}

//распаковка данных
func UnzipData(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	z, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
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

//читает и парсит json из файла
func JsonFromFile(file string) (map[string]interface{}, error) {
	js := make(map[string]interface{})
	//file := apppath + "/" + os.Args[1]
	//file := os.Args[1]
	data, err := FileRead(file)
	if err != nil {
		return nil, err
	}
	//LogPrintErrAndExit("JsonFromFile Error1: can't read file: "+file+"\n\n", err)

	js, err = FromJson(data)
	if err != nil {
		return nil, err
	}
	//LogPrintErrAndExit("JsonFromFile Error2: Unmarshal json error: "+file+"\n\n", err)

	return js, nil
}
