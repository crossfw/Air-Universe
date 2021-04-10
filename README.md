# Air-Universe

> 反馈 TG 群: https://t.me/Air_Universe <br>
> 此后端保留Xray完整配置, 并且不支持限速<br>
> 小白请使用 [XrayR后端](https://github.com/XrayR-project/XrayR)

## Features
- 支持3端(Shadowsocks, V2ray(VMess), Trojan) 单端口多用户
- Shadowsocks 单端口多用户 无须协议和混淆插件支持, 使用 AEAD 加密单端口 (原版 Clash 可用)
- V2ray(VMess) 支持 TCP 和 Websocket 可配合 TLS 传输, 证书可自定义(一键脚本不含此功能)也可自动生成
- Trojan 支持 TCP+TLS 
- 支持单启动多开节点，单服务器多节点无须多开程序，多个入站配合多节点ID, 流量分开统计
- 支持记录用户 IP, 但目前不可限制用户 IP 连接数
- 支持上报服务器负载和开机时间  
- 支持限速(请自行替换Xray内核 无教程) [点此下载](https://github.com/crossfw/Air-Universe/raw/master/scripts/proxy-core/xray_speedlimit)
- 审计规则默认屏蔽BT和内网IP, 可自行添加, 不支持从面板拉取
- 审计信息不会上报
- [一键脚本快速安装](https://github.com/crossfw/Air-Universe/tree/master/docs/TurnKey_cn.md)

## 文档
- 配置文件
  - [English](https://github.com/crossfw/Air-Universe/tree/master/docs/Doc_en.md)
  - [中文](https://github.com/crossfw/Air-Universe/tree/master/docs/Doc_cn.md)
- 一键脚本
  - [中文](https://github.com/crossfw/Air-Universe/tree/master/docs/TurnKey_cn.md)
- 手动安装
  - [中文](https://github.com/crossfw/Air-Universe/tree/master/docs/Install_cn.md)

## Thanks

* [Project X](https://github.com/XTLS/)
* [V2Fly](https://github.com/v2fly)
* [XrayR](https://github.com/XrayR-project/XrayR)
* [All stargazers](https://github.com/crossfw/Air-Universe/stargazers)

## Licence
 [GNU General Public License v3.0](https://github.com/crossfw/Air-Universe/blob/master/LICENSE)
## Stargazers over time

[![Stargazers over time](https://starchart.cc/crossfw/Air-Universe.svg)](https://starchart.cc/crossfw/Air-Universe)

