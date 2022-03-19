# Air-Universe 
## Introduction
Air-Universe is an open-source and free Middleware between SSPanel and V2Ray-core. It will be compatible with any version of v2ray-core(4.x).
## Features
* Sync users from your SSPanel to V2Ray-core
* Post traffic data to your SSPanel
* Fully customizable V2Ray-core profiles
* NO users count limit.

## TurnKey Install
```shell
wget -N --no-check-certificate --no-cache https://github.com/crossfw/Air-Universe/raw/master/scripts/xray_script/Install_server_with_xray.sh && bash Install_server_with_xray.sh
```

## Install in linux
1. Prepare V2Ray-core form [V2Ray Release](https://github.com/v2fly/v2ray-core/releases)
2. Make a v2ray config for you, you can refer from [Here](https://github.com/crossfw/Air-Universe/blob/master/example/v2ray-core_json/Single.json) and [V2Ray document](https://www.v2ray.com/) <br>
You must remain V2Ray-core API and a route rule for API, API port can change if you want, the default port for V2Ray-core API is 10085.
   
3. Download Air-Universe from [Air-Universe Release](https://github.com/crossfw/Air-Universe/releases)
4. Make a Air-Universe config from [Here](https://github.com/crossfw/Air-Universe/blob/master/example/v2rayssp_json/example.json) <br>
5. Start V2Ray-core first
```shell
./v2ray -c your_v2ray.json
```
6. After v2ray-core starts successfully, launch Air-Universe
```shell
./Air-Universe -C your_Air-Universe.json
```
7. Test connection.

## Air-Universe Configuration File Format

### Overview
Configuration of Air-Universe is a file with the following format. It includes "panel", "proxy", "sync"
```json
{
  "panel": {},
  "proxy": {},
  "sync": {}
}
```

> `panel`: [PanelObject](#panelobject)
 
 The settings about panel

> `proxy`: [ProxyObject](#proxyobject)

 Control local or remote Proxy-core

> `sync`: [SyncObject](#syncobject)

 Sync and retry settings


### PanelObject

`PanelObject` configuration format.
```json
{
  "type": "type of panel",
  "url": "https://SSPanel.address",
  "key": "SSPanel-Key",
  "node_ids": [1, 2]
}
```

> `type`: string

Which panel you are using. Now it supports 
- `sspanel` - [SSPanel-Uim](https://github.com/Anankke/SSPanel-Uim)

> `url`: string

Your SSPanel url. Make sure it start at "http" or "https".

> `key`: string

Your SSPanel's mu_key. Check from website root `./config/.config.php` if your panel is sspanel.

> `node_ids`: [uint32]

An array, each element of which is a nodeId you want to service. Please make it length equal to `proxy.in_tags`.<br>
The first element of `proxy.in_tags` will get users from the first element of `panel.node_ids`. Map by their index.

### ProxyObject
`ProxyObject` configuration format.
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
  "log_path": "./v2.log",
  "enable_sniffing": true,
  "cert": {
    "cert_path": "/path/to/certificate.crt",
    "key_path": "/path/to/key.key"
  }
}
```

> `type`: string

Which proxy you are using. Now it supports
- `v2ray` - [V2Ray-core](https://github.com/v2fly/v2ray-core)

> `alter_id`: string

V2Ray AlterId you want to set for every user. Please make sure it is equal to your panel setting.

> `auto_generate`: bool

If it's true, the Inbound to Xray will generate automatically.

> `in_tags`: [string]

An array includes Tags which v2ray inbound you want to add users.Please make it length equal to `panel.node_ids`.<br>
The first element of `proxy.in_tags` will get users from the first element of `panel.node_ids`. Map by their index.

> `api_address`: string

V2Ray-core api address, normally it will be "127.0.0.1"

> `api_port`: uint32

V2Ray-core api port, normally it will be "10085". You can change it in your v2ray-core config json.

> `log_path`: string

V2Ray-core log path, it will use to record users' IP address.

> `enable_sniffing`: true | false

When set to `false` Air-Universe will not set sniffing to all nodes.

> `cert`: cert object

Include domain certificate for tls

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

Interval time(second) between two synchronization.

> `fail_delay`: uint32

Retry delay time(second) if synchronization failure.

> `timeout`: uint32

HTTP connection request timeout(for connect to the panel)
