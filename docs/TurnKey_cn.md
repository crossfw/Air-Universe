# TurnKey Install
自动安装 Air-Universe + Xray, 在debian8测试通过, 不支持Centos
```shell
wget -N --no-check-certificate --no-cache https://github.com/crossfw/Air-Universe/raw/master/scripts/xray_script/Install_server_with_xray.sh && bash Install_server_with_xray.sh
```
- 仅适配[SSPanel-UIM](https://github.com/Anankke/SSPanel-Uim)
- 使用Xray做Proxy-core。
- 自动创建入站规则。
- 适配多入站，比如一个带ProxyProtocol的中转 和 一个直连入站，均采用面板节点ID， 流量分开统计。

实验性限速脚本
```shell
wget -N --no-check-certificate --no-cache https://github.com/crossfw/Air-Universe/raw/master/scripts/v2ray_script/Install_server_with_v2ray.sh && bash Install_server_with_v2ray.sh
```


## 详细说明
首先需要在面板上添加节点
请参考教程，比如[这个](https://soga.vaxilu.com/soga-v2ray/sspanel-v2ray) (逃
唯一的区别是，如果需要在中转后获取真实ip，请在v2ray(Vmess)或trojan节点地址配置时在最后加上"|relay=true"(不含引号)，
亦或者在节点类型中选择ss中转或v2ray中转。至此，3种协议均可获取中转真实IP， 不过要注意的是，在开启此功能时，TCP包开头必须携带ProxyProtocol，否则不予建立连接，
所以如果既要中转又要直连的，需要开2个不同的端口(配合不同节点ID)。<br>

注意 Shadowsocks AEAD单端口多用户需要更新面板到
[232c87c](https://github.com/Anankke/SSPanel-Uim/commit/232c87c0ff80d0118249d9c0eb161f869e7f4c5d)
之后, 且需将单端口承载用户的协议和混淆设为"origin"和"plain"(!注意,这个操作会使现有ssr单端口节点失效,谨慎操作!)<br>

如果使用自动生成证书的TLS, 请在节点信息后添加"|verify_cert=false"来跳过用户侧证书验证(需客户端支持)

### 需要输入的内容
```shell
########Air-Universe config#######\n
Enter node_ids, (eg 1,2,3): 2,3
Enter sspanel domain(https://): 1.1.1.1
Enter panel token: 123

```
- 节点ID列表, 不同id用英文逗号","分隔,最后一位不用加
- 面板地址(都2021年了还没HTTPS?) 输入域名即可,必须是https协议,否则你要自己去改配置文件.
- 面板密码

### 这个脚本会做什么
- 下载2个主程序到/usr/bin/au/
    - Air-Universe
    - Xray

- 配置文件目录/etc/au/
- 日志文件
    - Air-Universe log: /var/log/au.log
    - Xray log:/var/log/xr.log
    - Air-Universe日志文件每天6点清空
    - Xray日志文件每60s清空(用于统计ip)
    
### 脚本命令
```shell
$ /usr/bin/au/run.sh start|restart|stop
```
对应启动，重启，停止，默认会添加一条crontab保活并开机自启
