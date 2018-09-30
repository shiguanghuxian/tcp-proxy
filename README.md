# tcp-proxy
tcp代理程序，可以将服务器内网tcp服务代理到外网访问。

一般用于服务端内部服务(例如：redis、mysql、tabbitmq等)代理到公网方便查看调试自己开发的程序。

安装：

`go get github.com/shiguanghuxian/tcp-proxy`

配置：

```
bin/config/cfg.yaml

proxys: 
  - 
    name: test1 # 代理组名
    listen: 0.0.0.0:9898 # 代理组监听地址和端口
    reverse: # 被代理地址列表，可多个
      - 127.0.0.1:8114
      - 127.0.0.1:8115
```
