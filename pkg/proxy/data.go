package proxy

import "strings"

const data = `198.50.163.192:3129
170.79.168.225:8080
187.86.158.117:3128
213.174.89.7:3128
3.129.1.89:3128
217.114.10.92:8080
188.225.177.82:8080
186.118.169.85:8080
169.46.84.217:3128
91.195.205.39:8080
193.233.71.69:8080
13.76.195.117:8080
35.184.120.46:3128`

var ProxyList []Proxy

func init() {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		ProxyList = append(ProxyList, Proxy{parts[0], parts[1]})
	}
}
