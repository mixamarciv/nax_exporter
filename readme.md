ping exporter for prometheus.
```
git clone https://github.com/mixamarciv/nax_exporter.git
cd nax_exporter
export GOPATH=$(pwd)
go get -d
go build
./nax_exporter cfg.json
```


```
wget http://localhost:1001/ping?host=192.168.0.105&host=192.168.0.1&host=vrashke.net&host=google.com&host=ya.ru -O- -q
#worker 1
ping_ms{ip="ya.ru"} 30
ping_ttl{ip="ya.ru"} 55
ping_success{ip="ya.ru"} 1


#worker 4
ping_ms{ip="192.168.0.1"} 0.5
ping_ttl{ip="192.168.0.1"} 64
ping_success{ip="192.168.0.1"} 1


#worker 2
ping_ms{ip="192.168.0.105"} 1
ping_ttl{ip="192.168.0.105"} 128
ping_success{ip="192.168.0.105"} 1


#worker 3
ping_ms{ip="google.com"} 30
ping_ttl{ip="google.com"} 55
ping_success{ip="google.com"} 1


#worker 1
ping_ms{ip="vrashke.net"} 50
ping_ttl{ip="vrashke.net"} 54
ping_success{ip="vrashke.net"} 1
```


prometheus.yml
```
scrape_configs:
  - job_name: 'test123'
    scrape_interval: 10s
    metrics_path: '/ping'
    scheme: 'http'
    params:
        host: ['vrashke.net','anykeyadmin.info','192.168.0.105','192.168.0.30','192.168.0.31','192.168.0.105']
    static_configs:
      - targets: ['192.168.0.31:1000','192.168.0.30:3000'] # hosts where you start nax_exporter
```

