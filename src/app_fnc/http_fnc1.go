package app_fnc

import (
	"io/ioutil"
	"net/http"
)

//отправляем простой http запрос
func SendHttpRequest(urlStr string) (map[string]interface{}, error) {
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

	ret := make(map[string]interface{})
	ret["body"] = string(body)
	ret["response"] = resp

	return ret, nil
}
