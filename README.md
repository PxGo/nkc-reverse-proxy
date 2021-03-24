# NKC proxy go

golang实现的可配置的反向代理服务器

## usage

```bash
go run main.go path/to/config.json
```
## example

配置文件示例

[proxy.config.json](https://github.com/kccd/nkc-proxy-go/blob/main/proxy.config.json)
## benchmark

分别使用两个版本代理，在浏览器端for循环使用nkcAPI发起请求获得不同次数下的耗时对比，每个结果取自5次相同测试耗时最少的那个值

|proxy          | 10 times  | 50 times | 100 times |
|:--------------|:---------:|:--------:|:---------:|
|node version   | 7.46s     | 35.49s   | 71.03s    |
|golang version | 5.88s     | 19.28s   | 38.06s    |

golang版本http请求处理效率相较于node提升为 21% 到 46%

测试代码

```javascript
var arr = [];
console.time("cost")
for(var i=0;i<10;i++) {
  arr.push(nkcAPI("https://localhost/t/354150", "GET").then(() => {
    console.log("完成")
  }))
}
Promise.all(arr).then(() => {
  console.log("全部完成")
  console.timeEnd("cost")
})
```
