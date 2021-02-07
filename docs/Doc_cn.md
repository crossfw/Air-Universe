# V2ray-ssp
> This document was translated from English version by [DeepL](https://www.deepl.com/)
## 介绍
V2ray-ssp 是一个介于 SSPanel 和 V2ray-core 之间的开源和免费的中间件。它将兼容V2ray-core(4.x)的任何版本。
## 特征
* 将用户从SSPanel同步到V2ray-core。
* 将流量数据发布到SSPanel上
* 完全可定制的V2ray-核心配置文件。
* 无用户数限制。
## 在Linux上安装
1. 准备V2ray-核心表格[V2ray发布](https://github.com/v2fly/v2ray-core/releases)
2. 给自己做一个v2ray的配置，可以参考[这里](https://github.com/crossfw/V2ray-ssp/blob/master/example/v2ray-core_json/Single.json) 和 [V2ray文档](https://www.v2ray.com/) <br>
   你必须保留V2ray-core API和API的路由规则，API端口可以随意更改，V2ray-core API的默认端口是10085。
3. 从[V2ray-ssp Release](https://github.com/crossfw/V2ray-ssp/releases) 下载V2ray-ssp。
4. 从[这里](https://github.com/crossfw/V2ray-ssp/blob/master/example/v2rayssp_json/example.json) 进行V2ray-ssp配置
5. 先启动V2ray-core
```shell
./v2ray -c your_v2ray.json
```
6. 在V2ray-core 启动后再启动v2ray-ssp
```shell
./v2ray-ssp -C your_v2ray-ssp.json
```
7. 测试连接.

## V2ray-ssp 配置文件解析
```json
{
  "url": "https://SSPanel.address",
  "key": "SSPanel-Key",
  "node_id": 24,
  "alert_id": 1,
  "in_tags": ["p0"],
  "api_address": "127.0.0.1",
  "api_port": 10085,
  "sync_interval": 60,
  "fail_delay": 3
}

```

- url
    - 您的SSPanel的网址，请确保它以 "http "或 "https "开头。确保它以 "http "或 "https "开头。

- key
    - 你的SSPanel的mu_key。从SSPanel网站根目录下的./config/.config.php中查看。

- node_id
    - 您要建立的节点

- alert_id
    - 确保 alertId 与 SSPanel 的节点配置相同，错误的值将导致连接失败或内存泄漏。

- in_tags
    - 一个数组包括您想添加用户的v2ray-core inbound标签。

- api_address
    - V2ray-core的api地址，通常是 "127.0.0.1"

- api_port
    - V2ray-core的api端口，通常是 "10085"。你可以在v2ray-core配置的json中更改它。

- sync_interval
    - 两次同步的间隔时间（秒）。

- fail_delay
    - 同步失败时重试延迟时间(秒)。
    