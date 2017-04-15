package app_fnc

import (
	"strings"
)

func JsonToLines(js map[string]interface{}, prefix string) string {
	s := ""
	for key, val := range js {
		ikey := Sprintf("%v", key)
		if prefix != "" {
			ikey = prefix + "_" + ikey
		}
		switch tval := val.(type) {
		case bool, string, int, float64:
			s += ikey + " " + jsonToLines_getInterfaceVal(tval) + "\n"
		case []interface{}:
			s += interfaceArrToLines(tval, ikey) + "\n"
		case map[string]interface{}:
			s += JsonToLines(tval, ikey) + "\n"
		default:
			s += "#" + ikey + " cant get type info from val: " + interfaceToComment(tval) + "\n"
		}
	}
	return s
}

func interfaceArrToLines(jsarr []interface{}, prefix string) string {
	s := ""
	for key, val := range jsarr {
		ikey := Sprintf("%v", key)
		if prefix != "" {
			ikey = prefix + "_" + ikey
		}
		switch tval := val.(type) {
		case bool, string, int, float64:
			s += ikey + " " + jsonToLines_getInterfaceVal(tval) + "\n"
		case []interface{}:
			s += interfaceArrToLines(tval, ikey) + "\n"
		case map[string]interface{}:
			s += JsonToLines(tval, ikey) + "\n"
		default:
			s += "#" + ikey + " cant get type info from val: " + interfaceToComment(tval) + "\n"
		}
	}
	return s
}

func jsonToLines_getInterfaceVal(val interface{}) string {
	switch tval := val.(type) {
	case int, float64:
		return Sprintf("%v", tval)

	case string:
		if tval == "" {
			return "0"
		}
		ival := StrTrim(StrRegexpReplace(tval, "[^\\d\\.\\-e]*", ""))
		if ival == "" || ival == "." || strings.Index(ival, "..") >= 0 || strings.Count(ival, "e") > 1 || strings.Count(ival, "-") > 1 {
			return "0"
		}
		return ival

	case bool:
		if tval == false {
			return "0"
		}
		return "1"

	default:
		return "0\n#======> ERROR3000 !!!  cant get type info from val ERROR3000!!: " + interfaceToComment(tval)
	}
}

func interfaceToComment(js interface{}) string {
	s := Sprintf("%#v", js)
	s = strings.Replace(s, "\n", "\n#", -1)
	return s
}
