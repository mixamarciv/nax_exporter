package httpgetparsefuncs

import (
	a "app_fnc"

	//"strings"
)

//http://ethpool.org/api/miner_new/1d6604ffa0307db4df833cba721ce471e26f03cb
//{"address":"1d6604ffa0307db4df833cba721ce471e26f03cb","hashRate":"174.7 MH/s","reportedHashRate":"173.5 MH/s","blocks":[],
//"workers":
//{"rig1":{"worker":"rig1","hashrate":"81.1 MH/s","validShares":71,"staleShares":3,"invalidShares":0,"workerLastSubmitTime":1489936671,"invalidShareRatio":0,"reportedHashRate":"83.9 MH/s"},
//"rig2":{"worker":"rig2","hashrate":"93.7 MH/s","validShares":83,"staleShares":2,"invalidShares":0,"workerLastSubmitTime":1489936691,"invalidShareRatio":0,"reportedHashRate":"89.6 MH/s"}},
//"settings":{"miner":"1d6604ffa0307db4df833cba721ce471e26f03cb","email":"***amarciv@gmail.com","monitor":1,"name":null,"vote":0,"voteip":"","ip":"*.*.*.19"},
//"ethPerMin":0.000031958394923675686,"usdPerMin":0.0012320632369370375,"btcPerMin":0.0000012269594812677263,"avgHashrate":21452546.296296295,
//"credits":[{"miner":"1d6604ffa0307db4df833cba721ce471e26f03cb","credit":2193800000000,"time":"2017-03-19T15:18:37.000Z","balance":55909520963018430,"maxCredit":213061468005303.16}],
//"totalShareStats":{"valid":154,"invalid":0,"stale":5}}
func parse_balance_ethpool(data string) string {
	js := a.FromJsonStr([]byte(data))
	if _, b := js["error"]; b == true {
		return "ERROR_1_json_parse 1\n#" + data
	}
	//return data + "\n\n\n" +
	return a.JsonToLines(js, "")

}
