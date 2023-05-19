NKC-Reverse-Proxy
======
NKC-Reverse-Proxy is a powerful and efficient cross-platform reverse proxy service written in Golang. 

It now includes automated Let's Encrypt certificate retrieval.

[![Go Version](https://img.shields.io/badge/Go-v1.16-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/kccd/nkc-reverse-proxy)](https://goreportcard.com/report/github.com/kccd/nkc-reverse-proxy)
[![Downloads](https://img.shields.io/github/downloads/kccd/nkc-reverse-proxy/total)](https://github.com/kccd/nkc-reverse-proxy/releases)
[![References](https://img.shields.io/github/forks/kccd/nkc-reverse-proxy?label=references)](https://github.com/kccd/nkc-reverse-proxy/network/members)
[![License](https://img.shields.io/github/license/kccd/nkc-reverse-proxy)](https://github.com/kccd/nkc-reverse-proxy/blob/main/LICENSE)


Usage
-----

### Installation

You can install NKC-Reverse-Proxy by compiling the source code or downloading the pre-compiled binary.


#### Compile from source code

1.  `git clone https://github.com/kccd/nkc-reverse-proxy`
2.  `cd nkc-reverse-proxy`
3.  `go build .`

#### Download pre-compiled binary

Download the binary for your OS from the [releases page](https://github.com/kccd/nkc-reverse-proxy/releases).

### Running

To run NKC-Reverse-Proxy, use the following command in your terminal:

```bash
nkc-reverse-proxy -f /path/to/config/file
```
Please note that the name of the executable file may vary depending on the platform you are using. Make sure to use the correct name of the executable file for your platform when running the program.

Replace /path/to/config/file with the path to your configuration file. This command will start the reverse proxy server and load the configuration file specified.

### Configuration

#### console
Used to control the display of console logs.

yaml
```yaml
console:
  debug: false # Whether to display debug logs
  warning: false # Whether to display warning logs
  error: false # Whether to display error logs
  info: false # Whether to display info logs
```


#### proxy & maxIpCount

When the program is behind another proxy program, this value must be set to obtain the real IP and Port of the client.

yaml
```yaml
proxy: false # Whether to process after other proxy programs
maxIpCount: 1 # Maximum number of allowed IPs
```

For example, if the program is in the following scenario: Proxy1 -> Proxy2 -> Current program, to obtain the real IP and Port of the client, the following configuration needs to be done:

yaml
```yaml
proxy: true
maxIpCount: 2
```

#### req_limit
Access rate control.

yaml
```yaml
req_limit:
  - "50/10s 100" # Maximum of 50 requests processed in 10 seconds, with a cache of up to 100 requests, and no discrimination
  - "200/5m 400 ip" # Maximum of 200 requests processed in 5 minutes, with a cache of up to 400 requests, restricted by client IP
```
Valid time units are: `s` `m` `h` `d`.

For example, if the request limit is not exceeding 500 per minute and the maximum cache number is 2000, restricted by client IP:

yaml
```yaml
req_limit:
  - "500/1m 2000 ip"
```

#### servers
Used to configure the relevant information of the reverse proxy service.

For detailed configuration information, please refer to [config.yaml](https://github.com/kccd/nkc-reverse-proxy/blob/main/config.yaml)

Examples
-----

#### HTTP

yaml
```yaml
servers:
  -
    listen: 80
    name:     
      - www.example.com
    location: 
      -
        reg: "^\\/"
        pass:      
          - http://127.0.0.1:8080
        balance: random

```

#### HTTPS

yaml
```yaml
servers:
  -
    listen: 80
    name:
      - www.example.com
      - example.com
    location:
      -
        reg: "^\\/"
        redirect_code: 301
        redirect_url: "https://www.example.com"
  
  -
    listen: 443
    name: 
      - www.example.com
    ssl_cert: "/path/to/ssl/cert"
    ssl_key: "/path/to/ssl/key"
    location:
      -
        reg: "^\\/"
        pass:
          - http://127.0.0.1:9000
          - http://127.0.0.1:9001
          - http://127.0.0.1:9002
          - http://127.0.0.1:9003
        balance: random
```

#### SOCKET.IO

yaml
```yaml
servers:
  -
    listen: 80
    name:
      - www.example.com
      - example.com
    location:
      -
        reg: "^\\/"
        redirect_code: 301
        redirect_url: "https://www.example.com"
  
  -
    listen: 443
    name: 
      - www.example.com
    ssl_cert: "/path/to/ssl/cert"
    ssl_key: "/path/to/ssl/key" 
    location:
      -
        reg: "^\\/"
        pass:
          - http://127.0.0.1:9000
          - http://127.0.0.1:9001
          - http://127.0.0.1:9002
          - http://127.0.0.1:9003
        balance: random
      -
        reg: "^\\/socket.io\\/"
        pass:
          - http://127.0.0.1:12000
          - http://127.0.0.1:12001
          - http://127.0.0.1:12002
          - http://127.0.0.1:12003
        balance: ip_hash
```

License
-----
NKC-Reverse-Proxy is released under the [GNU Lesser General Public License v3.0](https://github.com/kccd/nkc-reverse-proxy/blob/main/LICENSE). This license grants you the freedom to use, modify, and distribute the software as you wish, subject to certain conditions.