# Air-Universe 手动安装说明

## 准备的程序和工具
- 机器一台，Windows 和 Linux 均可，架构请选择 release 中包含的架构，否则请自行编译。
- Xray 内核可执行文件一个，请自行选择合适的内核版本，配置 JSON 一个。
- Air-Universe 可执行文件一个，配置 JSON 一个。
- 机器需要有公网端口，否则外部无法连接。
- 一个安装正确的 SSPanel 或 V2board面板。

## 安装前准备
1. 确定所需的协议和端口
   

2. 前往 SSPanel 添加所需的的节点信息。可参照 [其他教程](https://soga.vaxilu.com/soga-v2ray/sspanel-v2ray) 添加.
    - 注意 Shadowsocks AEAD 单端口多用户需要更新面板到
      [232c87c](https://github.com/Anankke/SSPanel-Uim/commit/232c87c0ff80d0118249d9c0eb161f869e7f4c5d)
      之后, 且需将单端口承载用户的协议和混淆设为 "origin" 和 "plain" (注意,这个操作会使现有ssr单端口节点失效,谨慎操作!)
    - 如果使用自动生成证书的 TLS (VMess, Trojan), 请在节点信息后添加 "|verify_cert=false" 来跳过用户侧证书验证(需客户端支持)。
    - 如果需要在中转后获取真实ip，请在 V2Ray(VMess) 或 Trojan 节点地址配置时在最后加上 "|relay=true" (不含引号)，
      亦或者在节点类型中选择 Shadowsocks 中转或 V2Ray 中转。至此，3种协议均可获取中转真实IP， 不过要注意的是，在开启此功能时，TCP 包开头必须携带 ProxyProtocol, 否则不予建立连接，
    - 如果既要中转又要直连的，需要开2个不同的端口(配合不同节点ID)。暂时不支持一个端口同时接受中转和直连。
    

3. 程序下载地址
    - Air-Universe [下载地址](https://github.com/crossfw/Air-Universe/releases)
    - Xray-core [下载地址](https://github.com/XTLS/Xray-core/releases)
     

4. 配置文件制作。
   >代理内核配置文件请自行增加正确的 API 配置。
   > 
   > Air-Universe 默认自动配置入站。
    - Air-Universe
        - [文档地址](https://github.com/crossfw/Air-Universe/blob/master/docs/Doc_cn.md)
        - [样例地址](https://github.com/crossfw/Air-Universe/tree/master/configs/Air-Universe_json)
    - Xray
        - [文档地址](https://xtls.github.io/guide/document/)
        - [样例地址](https://github.com/crossfw/Air-Universe/blob/master/configs/xray_json/multiIn.json)
        
## 安装
1. 首先启动代理内核(Xray or V2Ray)
    ```shell
    $ /path2xray/xray -c /path2xrayConfig/config.json
    ```
    请自行选择后台运行方法(nohup 或 screen 等)
    

2. 运行 Air-Universe 
   请于代理内核启动后1-2秒后再运行
    ```shell
    $ /path2AirU/Air-Universe -c /path2AirUConfig/config.json
    ```
   请自行选择后台运行方法(nohup 或 screen 等)


3. 确认连接情况。

4. 自行选择保活及开机自启方案。

