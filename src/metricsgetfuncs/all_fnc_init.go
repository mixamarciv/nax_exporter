package metricsgetfuncs

var fnclist map[string](func(string, interface{}) string)

func Get_parser_fnc(fncname string) func(string, interface{}) string {
	f, b := fnclist[fncname]
	if !b {
		return func(string, interface{}) string {
			return "#error parser func " + fncname + " not found!\nERROR_nax_exporter{not_found_function=\"" + fncname + "\"} 1\n"
		}
	}
	return f
}

//сохраняем разобранные параметры запроса, что бы не парсить каждый раз
var cache_opt map[string]map[string]interface{}

func init() {
	cache_opt = make(map[string]map[string]interface{})

	fnclist = make(map[string](func(string, interface{}) string))

	fnclist["check_url"] = check_url
	fnclist["utilinfo"] = utilinfo
}
