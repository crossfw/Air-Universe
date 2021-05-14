# Air-Universe

开源多功能机场后端， 协议支持 V2Ray(VMess), Trojan, Shadowsocks(单端口多用户)；面板支持 SSPanel, v2board.

> 反馈 TG 群: https://t.me/Air_Universe <br>
> 喜欢本项目欢迎 Star⭐

## 文档

- [WIKI](https://github.com/crossfw/Air-Universe/wiki)
- [一键脚本](https://github.com/crossfw/Air-Universe/wiki/%E4%B8%80%E9%94%AE%E8%84%9A%E6%9C%AC%E5%AE%89%E8%A3%85)
- [Docker安装](https://github.com/crossfw/Air-Universe-DockerInstall)
- [配置文件文档](https://github.com/crossfw/Air-Universe/wiki/%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
- [手动安装说明](https://github.com/crossfw/Air-Universe/wiki/%E6%89%8B%E5%8A%A8%E5%AE%89%E8%A3%85)


## 功能介绍

> √: 使用该功能仅需使用一键脚本，无须用户干预<br>
> 手动: 使用该功能需要用户自行配置，需要用户手动修改 Xray 配置文件或更新 Xray 程序<br>
> ×: 目前不支持该功能


| 功能            | VMess | Trojan | Shadowsocks | Vless |
| --------------- | ----- | ------ | -----------| ------ |
| 自动配置节点    | √     | √      | √           | ×  |
| 获取用户信息    | √     | √      | √           | √ |
| 用户流量统计    | √     | √      | √           | √ |
| 服务器信息上报  | √     | √      | √           | √  |
| 自动申请tls证书 | √     | √      | √           | √ |
| 自动续签tls证书 | √     | √      | √           | √  |
| 在线人数统计    | √     | √      | √           | √ | 
| 在线用户限制    | ×     | ×      | ×           | × |
| 审计规则        | 手动     | 手动     | 手动     | 手动 |
| 审计上报        | ×     | ×          | ×       | × |
| 节点端口限速    | 手动     | 手动      | 手动      | 手动 | 
| 按照用户限速    | 手动     | 手动      | 手动      | 手动 |
| 自定义DNS       | 手动     | 手动     | 手动     | 手动 |

- 完全自定义 Xray 配置文件
- Shadowsocks 单端口多用户 无须协议和混淆插件支持, 使用 AEAD 加密单端口 (原版 Clash 可用)
- 支持单启动多开节点，单服务器多节点无须多开程序，多个入站配合多节点ID, 流量分开统计
- 审计规则默认屏蔽BT和内网IP, 可自行添加, 不支持从面板拉取
- 可自动生成配置的节点类型
    - Shadowsocks
    - VMess(V2ary) 
      - 传输方式: TCP, Websocket, KCP, HTTP
      - 传输层加密: TLS
    - Trojan (TCP+TLS)

## 支持前端

| 前端        | v2ray | trojan | shadowsocks |
| ----------- | ----- | ------ | ---------- |
| sspanel-uim | √     | √      | √  |
| v2board     | √     | √      | √          |


## Thanks

* [Project X](https://github.com/XTLS/)
* [V2Fly](https://github.com/v2fly)
* [XrayR](https://github.com/XrayR-project/XrayR) - 另一个开源后端, 原生支持限速限用户,但不支持自定义Xray配置.
* [All stargazers](https://github.com/crossfw/Air-Universe/stargazers)

## Licence

[GNU General Public License v3.0](https://github.com/crossfw/Air-Universe/blob/master/LICENSE)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/crossfw/Air-Universe.svg)](https://starchart.cc/crossfw/Air-Universe)

