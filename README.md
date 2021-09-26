# tcpsidecar

## Introduction

用 Go 实现的端口转发工具，分为 Data 和 Control 两种节点。

Data 节点有多个，负责每个程序的负责数据层面的工作

1. 监听端口
2. 转发数据
3. 接收控制信息

Control 节点全局只有一个，负责控制层面的工作

1. 监控所有 Data 状态
2. 控制 Data 节点的转发关系

## Documentation



## TODO

* [x] TCP 协议的转发器
* [ ] 流量统计
* [ ] RPC 控制转发器的转发规则
* [ ] 抽象为 Data 节点

