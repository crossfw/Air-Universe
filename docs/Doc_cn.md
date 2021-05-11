# Air-Universe
## 介绍
Air-Universe 是一个介于 SSPanel、V2board 和 V2ray-core 或 Xray-core 之间的开源和免费的中间件。它将兼容 V2ray-core(4.x) 或 Xray(1.3+) 的任何版本。
## 特征
- 支持3端(Shadowsocks, V2ray(Vmess), Trojan) 单端口多用户
- **Shadowsocks 单端口多用户 无须协议和混淆插件支持, 使用AEAD加密单端口**(原版 Clash 可用)
- V2ray(VMess) 支持 TCP 和 Websocket 可配合 TLS 传输, 证书可自定义(一键脚本不含此功能)也可自动生成
- Trojan 支持TCP+TLS
- 支持多个入站配合多节点ID, 流量分开统计
- 支持记录用户IP, 但目前不可限制
- 支持实验性限速(V2Ray后端)
- 审计规则默认屏蔽BT和内网IP, 可自行添加, 不支持从面板拉取
- 审计信息**不会**上报

## Air-Universe 配置文件解析
### Overview
配置文件包含 "log", "panel", "proxy", "sync"
```json
{
  "log": {},
  "panel": {},
  "proxy": {},
  "sync": {}
}
```
> `log`: [LogObject](#logobject)

日志的相关配置

> `panel`: [PanelObject](#panelobject)

前端Panel的配置

> `proxy`: [ProxyObject](#proxyobject)

Proxy-core(V2Ray 或 Xray)的API通信配置等

> `sync`: [SyncObject](#syncobject)

信息同步设置

### LogObject
`logObject` 格式。
```json
{
  "log_level": "info",
  "access": "/var/log/au/au.log"
}
```
> `log_level`: “debug” | “info” | “warning” | “error” | “panic”

日志的级别, 日志需要记录的信息. 默认值为 "info"。

- "debug"：调试程序时用到的输出信息。同时包含所有 "info" 内容。
- "info"：运行时的状态信息等，不影响正常使用。同时包含所有 "warning" 内容。
- "warning"：发生了一些并不影响正常运行的问题时输出的信息，但有可能影响用户的体验。同时包含所有 "error" 内容。
- "error"：遇到了无法正常运行的问题，但是不影响主体运行。
- "panic"：程序崩溃前的日志，需要立即处理。

> `access`: string

日志保存路径

### PanelObject

`PanelObject` 格式。
```json
{
  "type": "type of panel",
  "url": "https://SSPanel.address",
  "key": "SSPanel-Key",
  "node_ids": [1, 2],
  "nodes_type": ["vmess", "ss"],
  "nodes_proxy_protocol": [false, true]
}
```

> `type`: string

您使用的面板类型， 目前支持
- `sspanel` - [SSPanel-Uim](https://github.com/Anankke/SSPanel-Uim)
- `v2board` - [v2board](https://github.com/v2board/v2board)

> `url`: string

面板的通信地址，确保使用 "http://" 或 "https://" 开始

> `key`: string

面板通信密钥
> `node_ids`: [uint32]

一个数组，每个元素都是你将要部署的nodeId，请让它的长度等于`proxy.in_tags`。
`proxy.in_tags`的第一个元素将从`panel.node_ids`的第一个元素中获取用户。按其索引进行映射。

> `nodes_type` [string]

节点类型，仅在使用v2board时使用，多个节点请依次输入，支持类型
- `vmess` - V2Ray(VMess)
- `trojan` - Trojan
- `ss` - Shadowsocks

> `nodes_proxy_protocol` [bool]

节点是否接收ProxyProtocol(中转后获取真实IP)，仅在使用v2board时使用，输入为 bool 数组


### ProxyObject
`ProxyObject` 配置格式
```json
{
  "type": "type of proxy",
  "alter_id": 1,
  "auto_generate": true,
  "in_tags": [
    "p0",
    "p1"
  ],
  "api_address": "127.0.0.1",
  "api_port": 10085,
  "force_close_tls": false,
  "log_path": "./v2.log",
  "cert": {
    "cert_path": "/path/to/certificate.crt",
    "key_path": "/path/to/key.key"
  },
  "speed_limit_level": [0, 2, 10, 30, 60, 100, 150, 250, 400]
}
```

> `type`: string

您使用的代理内核，目前支持
- `v2ray` - [V2Ray-core](https://github.com/v2fly/v2ray-core)
- `xray` - [Xray](https://github.com/XTLS/Xray-core)

> `alter_id`: string

手动配置入站代理时分配给 VMess 用户的alter_id。

> `auto_generate`: bool

是否自动生成入站代理（仅适配Xray）
> `in_tags`: [string]

数组中包含了您要添加用户的 Inbound Tag，请使其长度等于`panel.node_ids`。
`proxy.in_tags`的第一个元素将从`panel.node_ids`的第一个元素中获取用户。按其索引进行映射。

> `api_address`: string

代理内核 API 地址
> `api_port`: uint32

代理内核 API 端口

> `force_close_tls`: bool

强制不启用 tls，即使面板节点配置了 tls，用于 ws 反代的情况. 

> `log_path`: string

代理内核的日志输出地址（用于统计用户在线ip）
> `cert`: cert object

TLS 所需的证书，默认路径 `/usr/local/share/server.crt` 与 `/usr/local/share/server.key`，仅在启用TLS时需要

> `speed_limit_level`: [float32]

V2Ray 等级对应的限速配置，确保0级为不限速，如配置流量耗尽限速1M，请配置等级1为 1M， 请按限速从小到大配置。单位Mbps
### SyncObject
`SyncObject` configuration format.
```json
{
  "interval": 60,
  "fail_delay": 5,
  "timeout": 5,
  "post_ip_interval": 90
}
```
> `interval`: uint32

于前端面板的同步时间间隔（秒）。
> `fail_delay`: uint32

同步失败后重试等待时间（秒）(目前版本已弃用，默认超时10s)。
> `timeout`: uint32

HTTP(S) 请求超时时间（秒）    

> `post_ip_interval`: uint32

同步 alive IP 的时间

## 最小配置文件启动
```json
{
  "panel": {
    "url": "https://SSPanel.address",
    "key": "SSPanel-Key",
    "node_ids": [1]
  }
}
```
此时完整的配置（默认配置为）
```json
{
  "panel": {
    "type": "sspanel",
    "url": "https://SSPanel.address",
    "key": "SSPanel-Key",
    "node_ids": [1]
  },
  "proxy": {
    "type": "xray",
    "alter_id": 1,
    "auto_generate": true,
    "force_close_tls": false,
    "in_tags": [
      "p0"
    ],
    "api_address": "127.0.0.1",
    "api_port": 10085,
    "log_path": "./v2.log",
    "speed_limit_level": [0, 0.2, 3, 7, 13, 19, 25, 38, 63]
  },
  "sync": {
    "interval": 60,
    "fail_delay": 5,
    "timeout": 10,
    "post_ip_interval": 90
  }
}
```
需要说明的是，proxy.in_tags 会根据 node_ids 自动补全。补全策略请参看 ./cmd/Air-Universe/cfg.go 文件。