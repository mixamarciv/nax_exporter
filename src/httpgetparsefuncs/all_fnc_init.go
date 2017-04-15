package httpgetparsefuncs

var fnclist map[string](func(string) string)

func Get_parser_fnc(fncname string) func(string) string {
	f, b := fnclist[fncname]
	if !b {
		return func(string) string {
			return "#error parser func " + fncname + " not found!\nERROR_nax_exporter{not_found_function=\"" + fncname + "\"} 1"
		}
	}
	return f
}

func init() {
	fnclist = make(map[string](func(string) string))

	fnclist["jsondata2lines4prometheus"] = parse_jsondata2lines4prometheus
	fnclist["ethminer"] = parse_ethminer
	fnclist["balance_coinminepl"] = parse_balance_coinminepl
	fnclist["balance_ethpool"] = parse_balance_ethpool
}
