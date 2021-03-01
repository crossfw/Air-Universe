# Air-Universe
## 介绍
Air-Universe 是一个介于 SSPanel 和 V2ray-core 或 Xray-core 之间的开源和免费的中间件。它将兼容 V2ray-core(4.x) 或 Xray(1.3+) 的任何版本。
## 特征
- 支持3端(Shadowsocks, V2ray(Vmess), Trojan) 单端口多用户
- **Shadowsocks 单端口多用户 无须协议和混淆插件支持, 使用AEAD加密单端口**(原版 Clash 可用)
- V2ray(VMess) 支持 TCP 和 Websocket 可配合 TLS 传输, 证书可自定义(一键脚本不含此功能)也可自动生成
- Trojan 支持TCP+TLS
- 支持多个入站配合多节点ID, 流量分开统计
- 支持记录用户IP, 但目前不可限制
- 不支持限速
- 审计规则默认屏蔽BT和内网IP, 可自行添加, 不支持从面板拉取
- 审计信息**不会**上报

## Air-Universe 配置文件解析
### Overview
配置文件包含 "panel", "proxy", "sync"
```json
{
  "panel": {},
  "proxy": {},
  "sync": {}
}
```

> `panel`: [PanelObject](#panelobject)

前端Panel的配置

> `proxy`: [ProxyObject](#proxyobject)

Proxy-core(V2Ray 或 Xray)的API通信配置等

> `sync`: [SyncObject](#syncobject)

信息同步设置

### PanelObject

`PanelObject` 格式。
```json
{
  "type": "type of panel",
  "url": "https://SSPanel.address",
  "key": "SSPanel-Key",
  "node_ids": [1, 2]
}
```

> `type`: string

您使用的面板类型， 目前支持
- `sspanel` - [SSPanel-Uim](https://github.com/Anankke/SSPanel-Uim)

> `url`: string

面板的通信地址，确保使用 "http://" 或 "https://" 开始

> `key`: string

面板通信密钥
> `node_ids`: [uint32]

一个数组，每个元素都是你将要部署的nodeId，请让它的长度等于`proxy.in_tags`。
`proxy.in_tags`的第一个元素将从`panel.node_ids`的第一个元素中获取用户。按其索引进行映射。

### ProxyObject
`ProxyObject` 配置格式
```json
{
  "type": "type of proxy",
  "alert_id": 1,
  "auto_generate": true,
  "in_tags": [
    "p0",
    "p1"
  ],
  "api_address": "127.0.0.1",
  "api_port": 10085,
  "log_path": "./v2.log",
  "cert": {
    "cert_path": "/path/to/certificate.crt",
    "key_path": "/path/to/key.key"
  }
}
```

> `type`: string

您使用的代理内核，目前支持
- `v2ray` - [V2Ray-core](https://github.com/v2fly/v2ray-core)
- `xray` - [Xray](https://github.com/XTLS/Xray-core)

> `alert_id`: string

手动配置入站代理时分配给 VMess 用户的alert_id。

> `auto_generate`: bool

是否自动生成入站代理（仅适配Xray）
> `in_tags`: [string]

数组中包含了您要添加用户的 Inbound Tag，请使其长度等于`panel.node_ids`。
`proxy.in_tags`的第一个元素将从`panel.node_ids`的第一个元素中获取用户。按其索引进行映射。

> `api_address`: string

代理内核 API 地址
> `api_port`: uint32

代理内核 API 端口
> `log_path`: string

代理内核的日志输出地址（用于统计用户在线ip）
> `cert`: cert object

TLS 所需的证书，如果留空则会自动生成（仅Xray），请在客户端跳过证书检查。
### SyncObject
`SyncObject` configuration format.
```json
{
  "interval": 60,
  "fail_delay": 5,
  "timeout": 5
}
```
> `interval`: uint32

于前端面板的同步时间间隔（秒）。
> `fail_delay`: uint32

同步失败后重试等待时间（秒）。
> `timeout`: uint32

HTTP(S) 请求超时时间（秒）    

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
    "alert_id": 1,
    "auto_generate": true,
    "in_tags": [
      "p0"
    ],
    "api_address": "127.0.0.1",
    "api_port": 10085,
    "log_path": "./v2.log",
  },
  "sync": {
    "interval": 60,
    "fail_delay": 5,
    "timeout": 5
  }
}
```
需要说明的是，proxy.in_tags 会根据 node_ids 自动补全。补全策略请参看 ./cmd/Air-Universe/cfg.go 文件。