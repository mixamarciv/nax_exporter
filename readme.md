ping exporter for prometheus.
```
git clone https://github.com/mixamarciv/nax_exporter.git
cd nax_exporter
export GOPATH=$(pwd)
go get -d
go build
./nax_exporter cfg.json
```


ping_192.168.0.31_ms 0.048
ping_192.168.0.31_ttl 64
ping_192.168.0.31_success 1
```
