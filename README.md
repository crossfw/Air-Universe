# Air-Universe
English document → [here](https://github.com/crossfw/Air-Universe/tree/master/docs/Doc_en.md) <br>
中文文档(咕咕咕) → [here](https://github.com/crossfw/Air-Universe/tree/master/docs/Doc_cn.md)

## Bugs:
- None
## TODO List:
- ~~Fix known bugs.~~
- ~~Auto generate Xray-core configuration from SSPanel.~~  Finished
- Support speed limit.
- ~~Record users ip.~~  Finished
- ~~Support Trojan protocol.~~ Finished
- ~~Turnkey installation script.~~ Finished
- ~~Support Shadowsocks multiuser in single port.~~ Finished

## TurnKey Install
自动安装 Air-Universe + Xray
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
所以如果既要中转又要直连的，需要开2个不同的端口(配合不同节点ID)。

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
    