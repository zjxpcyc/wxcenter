# 统一的微信服务中心

微信项目越来越多, 同时微信的各种服务也在各个项目中写死。这里抽离出来，实现微信服务的统一管理

服务重启不影响业务

依赖自己的包  
- https://github.com/zjxpcyc/wechat.v2
- https://github.com/zjxpcyc/tinylogger
- https://github.com/zjxpcyc/wechat-scheduler



## 启动
`go build` 之后
```bash
// -p 是端口, -v 显示版本号
./wxcenter -p=8080
```