package httpgetparsefuncs

import (
	a "app_fnc"
	"regexp"
	"strings"
)

//{"result": ["7.3", "199", "84138;633;0", "29854;25513;28770", "0;0;0", "off;off;off", "", "eth-eu.coinmine.pl:4000", "0;0;0;0"]}
//{"result": ["7.3", "199", "89637;741;0", "30193;29616;29827", "0;0;0", "off;off;off", "", "eth-eu.coinmine.pl:4000", "0;0;0;0"]}
/***********
hz_param1 7.3
hz_param2 218
hashrate_total 89664
hz_param3_2 806
hz_param3_3 0
hashrate_video_1 30179
hashrate_video_2 29623
hashrate_video_3 29860
temperature_1 0
temperature_2 0
temperature_3 0
************/
func parse_ethminer(data string) string {
	re := regexp.MustCompile("(?is)<body[^<]*<")
	s1 := re.FindString(data)
	s1 = a.StrRegexpReplace(s1, "[^\\{]*\\{", "{")
	s1 = a.StrRegexpReplace(s1, "[^\\}]*$", "")

	js := a.FromJsonStr([]byte(s1))
	if _, b := js["error"]; b == true {
		return "ERROR_1_json_parse 1"
	}

	result, b := js["result"]
	if b != true {
		return "ERROR_2_json_parse 1"
	}

	ar := result.([]interface{})

	ret := ""
	ret += "hz_param1 " + ar[0].(string) + "\n"
	ret += "hz_param2 " + ar[1].(string) + "\n"

	aa := strings.Split(ar[2].(string), ";")
	ret += "hashrate_total " + aa[0] + "\n"
	ret += "hz_param3_2 " + aa[1] + "\n"
	ret += "hz_param3_3 " + aa[2] + "\n"

	aa = strings.Split(ar[3].(string), ";")
	for i, v := range aa {
		ret += "hashrate_video_" + a.Itoa(i+1) + " " + v + "\n"
	}

	aa = strings.Split(ar[4].(string), ";")
	for i, v := range aa {
		ret += "temperature_" + a.Itoa(i+1) + " " + v + "\n"
	}

	return ret
}

//{"getuserbalance":{"version":"1.0.0","runtime":2.7709007263184,"data":{"confirmed":0.6283486,"unconfirmed":0.03604028,"orphaned":0}}}
//coinmine.pl
func parse_balance_coinminepl(data string) string {
	js := a.FromJsonStr([]byte(data))
	if _, b := js["error"]; b == true {
		return "ERROR_1_json_parse 1\n#" + data
	}

	js, b := js["getuserbalance"].(map[string]interface{})
	if b != true {
		return "ERROR_2_json_parse 1"
	}

	js, b = js["data"].(map[string]interface{})
	if b != true {
		return "ERROR_3_json_parse 1"
	}

	confirmed, b := js["confirmed"]
	if b != true {
		return "ERROR_4_json_parse 1"
	}

	ret := ""
	ret += "confirmed " + a.Sprintf("%f", confirmed) + "\n"

	unconfirmed, b := js["unconfirmed"]
	if b != true {
		return "ERROR_5_json_parse 1"
	}
	ret += "unconfirmed " + a.Sprintf("%f", unconfirmed) + "\n"

	return ret
}
