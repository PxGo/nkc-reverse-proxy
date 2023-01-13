# <center>nkc-reverse-proxy</center> 

#### <center>一个简单易用的反向代理服务</center>

---

## 目录
+ [安装](#install)
+ [配置说明](#configs)
  + [console](#configs_console)
  + [proxy](#configs_proxy)
  + [maxIpCount](#configs_proxy)
  + [req_limit](#configs_req_limit)
  + [servers](#configs_servers) 
    + [HTTP 网站配置](#configs_servers_http)
    + [HTTPS 网站配置](#configs_servers_https)
  
## <span id="install">安装</span>

```
$ git clone https://github.com/kccd/nkc-reverse-proxy.git
$ cd nkc-reverse-proxy
$ go build .
$ copy configs.template.yaml configs.yaml

# windows
$ ./nkc-reverse-proxy.exe
# linux
$ ./nkc-reverse-proxy

# 程序在启动时，会将项目根目录下的 configs.yaml 作为配置文件，你还可以手动指定配置文件所在位置，如：
$ ./nkc-reverse-proxy /workspace/proxy/configs.yaml
```

## <span id="configs">配置说明</span>

### <span id="configs_console">console</span>

用于控制控制台日志的显示。

```
console:
  debug: false # 是否显示 debug 日志  
  warning: false # 是否显示 warning 日志
  error: false # 是否显示 error 日志
  info: false # 是否显示 info 日志
```

### <span id="configs_proxy">proxy & maxIpCount</span>

当程序处于其他代理程序之后时需设置该值才能获取到客户端的真实 `IP `和 `Port`。

```
proxy: false # 是否处理其他代理程序之后
maxIpCount: 1 # 最大允许的 IP 数量
```
例如在 ` 代理1 -> 代理2 -> 当前程序 ` 这种情况下，想要获取客户端真实 `IP` 和 `Port`，就需要做如下配置：
```
proxy: true
maxIpCount: 2
```

### <span id="configs_req_limit">req_limit</span>

访问速率控制。

```
req_limit:
  - "50/s 100" # 每秒最多处理 50 个请求，缓存请求数不超过 100，无差别限制
  - "5/s 10 ip" # 每秒最多处理 5 个请求，缓存请求数不超过 10，根据客户端 IP 限制
```

合法的单位时间处理条数有：`num/s` `num/m` `num/h` `num/d`。

例如每分钟处理请求数不超过 `500` 且最大缓存数为 `2000`，根据客户端 `IP` 限制：
```
req_limit: 
  - "500/m 2000 ip"
```

### <span id="configs_servers">servers</span>

用于配置反向代理服务的相关信息。

```
servers:
  -
    listen: 443 # 服务暴露的端口
    name: # 允许链接的域
      - "127.0.0.1"
      - "localhost"
    ssl_key: "/ssl/test.key"  # SSL 证书文件路径
    ssl_cert: "/ssl/test.crt" # SSL 证书文件路径
    req_limit: # 访问速率限制
      - "50/s 500" 
      - "10/s 100 ip"
    location: # 根据路径匹配服务
      -    
        reg: "^\\/" # 请求路径正则
        pass: # 目标服务
          - "http://127.0.0.1:8080"
        balance: "random" # 负载均衡类型 random, ip_hash
        req_limit: # 访问速率限制
          - "10/s 50"
          - "5/s 30 ip"
      - 
        reg: "^\\/socket\\.io\\/"
        pass: 
          - "http://127.0.0.1:9090"
        balance: "ip_hash"
        req_limit:
          - "50/s 100"
          - "1/s 5 ip"    
      -
        reg: "^\\/old-home"
        redirect_code: 301 # 重定向状态码
        redirect_url: "https://127.0.0.1/home" # 重定向链接
```

以上就是 `servers` 中可能出现的配置选项，下面是一些例子。

#### <span id="configs_servers_http">1、HTTP 网站配置</span>

```
servers:
  - 
    listen: 80
    name: 
      - "www.domain.com"
    location:
      -
        reg: "^\\/"
        pass: 
          - "http://127.0.0.1:8080"
          - "http://127.0.0.1:8081"
          - "http://127.0.0.1:8081"
          - "http://127.0.0.1:8082"
        balance: "random"   
        req_limit: 
          - "50/s 300"
          - "3/s 10 ip"
```

#### <span id="configs_servers_https">2、HTTPS 网站配置</span>

```
servers:
  - 
    listen: 443
    name: 
      - "www.domain.com"
    ssl_key: "/ssl/www.domain.com.key"
    ssl_cert: "/ssl/www.domian.com.crt"
    location:
      -
        reg: "^\\/"
        pass: 
          - "http://127.0.0.1:8080"
          - "http://127.0.0.1:8081"
          - "http://127.0.0.1:8081"
          - "http://127.0.0.1:8082"
        balance: "random"   
        req_limit: 
          - "50/s 300"
          - "3/s 10 ip"
  -
    listen: 80
    name: 
      - "www.domain.com"
      - "domain.com"        
    location:
      -
        reg: "^\\/"
        redirect_code: 301
        redirect_url: "https://www.domain.com"  
```