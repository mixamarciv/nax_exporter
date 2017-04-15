package httpgetparsefuncs

import (
	a "app_fnc"
)

func parse_jsondata2lines4prometheus(data string) string {
	js := a.FromJsonStr([]byte(data))
	if _, b := js["error"]; b == true {
		return "ERROR_1_parse_jsondata2lines4prometheus 1\n#" + data
	}
	//return data + "\n\n\n" +
	return a.JsonToLines(js, "")
}
