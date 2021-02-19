# Air-Universe
English document → [here](https://github.com/crossfw/Air-Universe/tree/master/docs/Doc_en.md) <br>
中文文档(咕咕咕) → [here](https://github.com/crossfw/Air-Universe/tree/master/docs/Doc_cn.md)

>新人第一次写大项目，也是第一次写golang，请多多指教

## Bugs:
- None
## TODO List:
- ~~Fix known bugs.~~
- ~~Auto generate Xray-core configuration from SSPanel.~~  Finished
- ~~Record users ip.~~  Finished
- ~~Support Trojan protocol.~~ Finished
- ~~Turnkey installation script.~~ Finished
- ~~Support Shadowsocks multiuser in single port.~~ Finished
- Support speed limit.
- Limit users IP count.
- Support all platform turnkey script.

## 注意, 未经严格测试, 出问题了造成经济损失一概不负责
## Features
- 支持3端(Shadowsocks, V2ray, Trojan) 单端口多用户
- **Shadowsocks 单端口多用户 无须协议和混淆插件支持, 使用AEAD加密单端口**
- V2ray 支持 tcp和Websocket 可配合TLS传输, 证书可自定义(一键脚本不含此功能)也可自动生成
- Trojan 支持TCP+TLS 
- 支持多个入站配合多节点ID, 流量分开统计
- 支持记录用户IP, 但目前不可限制
- 不支持限速
- 审计规则默认屏蔽BT和内网IP, 可自行添加, 不支持从面板拉取
- 审计信息**不会**上报

## TurnKey Install
自动安装 Air-Universe + Xray, 在debian8测试通过, 不支持Centos
```shell
wget -N --no-check-certificate --no-cache https://github.com/crossfw/Air-Universe/raw/master/scripts/xray_script/Install_server_with_xray.sh && bash Install_server_with_xray.sh
```
- 仅适配[SSPanel-UIM](https://github.com/Anankke/SSPanel-Uim)
- 使用Xray做Proxy-core, 同时也支持V2Ray-Core
- 自动创建入站规则
- 适配多入站，比如一个带ProxyProtocol的中转 和 一个直连入站，均采用面板节点ID， 流量分开统计。

### 详细说明
首先需要在面板上添加节点
请参考教程，比如[这个](https://soga.vaxilu.com/soga-v2ray/sspanel-v2ray) (逃
唯一的区别是，如果需要在中转后获取真实ip，请在v2ray或trojan节点地址配置时在最后加上"|relay=true"(不含引号)， 
亦或者在节点类型中选择ss中转或v2ray中转。至此，3种协议均可获取中转真实IP， 不过要注意的是，在开启此功能时，TCP包开头必须携带ProxyProtocol，否则不予建立连接，
所以如果既要中转又要直连的，需要开2个不同的端口(配合不同节点ID)。<br>

注意 Shadowsocks AEAD单端口多用户需要更新面板到
[232c87c](https://github.com/Anankke/SSPanel-Uim/commit/232c87c0ff80d0118249d9c0eb161f869e7f4c5d)
之后, 且需将单端口承载用户的协议和混淆设为"origin"和"plain"(!注意,这个操作会使现有ssr单端口节点失效,谨慎操作!)<br>

如果使用自动生成证书的TLS, 请在节点信息后添加"|verify_cert=false"来跳过用户侧证书验证(需客户端支持)

#### 需要输入的内容
```shell
########Air-Universe config#######\n
Enter node_ids, (eg 1,2,3): 2,3
Enter sspanel domain(https://): 1.1.1.1
Enter panel token: 123

```
- 节点ID列表, 不同id用英文逗号","分隔,最后一位不用加
- 面板地址(都2021年了还没HTTPS?) 输入域名即可,必须是https协议,否则你要自己去改配置文件.
- 面板密码

#### 这个脚本会做什么
- 下载2个主程序到/usr/bin/au/
    - Air-Universe
    - Xray
    
- 配置文件目录/etc/au/
- 日志文件
    - Air-Universe log: /var/log/au.log
    - Xray log:/var/log/xr.log
    - Air-Universe日志文件每天6点清空
    - Xray日志文件每60s清空(用于统计ip)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/crossfw/Air-Universe.svg)](https://starchart.cc/crossfw/Air-Universe)
