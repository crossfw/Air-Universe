# TurnKey Install
Air-Universe 管理脚本(安装+拆卸+申请证书)
```shell
wget -N --no-check-certificate --no-cache https://github.com/crossfw/Air-Universe/raw/master/scripts/xray_script/Install_server_with_xray.sh && bash Install_server_with_xray.sh
```
升级Xray
```shell
bash -c "$(curl -L https://github.com/crossfw/Xray-install/raw/main/install-release.sh)" @ install
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

注意 证书务必使用 `fullchain.cer` 否则可能会导致无法连接问题。
### 需要输入的内容
```shell
########Air-Universe config#######
Enter node_ids, (eg 1,2,3): 1,2,1
Enter panel domain(Include https:// or http://): https://xxx.cloud
Enter panel token: xxxxxxx

Choose panel type:
  1. SSPanel
  2. V2board
Choose panel type: 2
Enter nodes type, (eg vmess,ss): "vmess","vmess","ss"

```
- 节点ID列表, 不同id用英文逗号","分隔,最后一位不用加
- 面板地址 请携带https:// 或 http://
- 面板密码
- 选择面板类型（若选择 v2board 则需要输入一下两项）
- 节点类型，请输入 "vmess", "trojan", "ss" 类型选项。 请和最开头的节点顺序一致。不要忘了加双引号。

### 这个脚本会做什么
- 下载2个主程序到/usr/local/bin/
    - Air-Universe
    - Xray

- Air-Universe 配置文件 /usr/local/etc/au/au.json
- Xray 配置文件 /usr/local/etc/xray/config.json
- 日志文件
    - Air-Universe log: /var/log/au/au.log
    - Xray log: /var/log/au/xr.log
    - Xray日志文件每300s清空(用于统计ip)
    
### 脚本命令
```shell
$ systemctl start|stop|restart au
```
对应启动，重启，停止, 默认开机自启
