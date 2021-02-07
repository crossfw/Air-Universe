# V2ray-ssp 
## Introduction
V2ray-ssp is an open-source and free Middleware between SSPanel and V2ray-core. It will be compatible with any version of v2ray-core(4.x).
## Features
* Sync users from your SSPanel to V2ray-core
* Post traffic data to your SSPanel
* Fully customizable V2ray-core profiles
* NO users count limit.
## Install in linux
1. Prepare V2ray-core form [V2ray Release](https://github.com/v2fly/v2ray-core/releases)
2. Make a v2ray config for you, you can refer from [Here](https://github.com/crossfw/V2ray-ssp/blob/master/example/v2ray-core_json/Single.json) and [V2ray document](https://www.v2ray.com/) <br>
You must remain V2ray-core API and a route rule for API, API port can change if you want, the default port for V2ray-core API is 10085.
   
3. Download V2ray-ssp from [V2ray-ssp Release](https://github.com/crossfw/V2ray-ssp/releases)
4. Make a V2ray-ssp config from [Here](https://github.com/crossfw/V2ray-ssp/blob/master/example/v2rayssp_json/example.json) <br>
5. Start V2ray-core first
```shell
./v2ray -c your_v2ray.json
```
6. After v2ray-core starts successfully, launch V2ray-ssp
```shell
./v2ray-ssp -C your_v2ray-ssp.json
```
7. Test connection.

## V2ray-ssp Config explain
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
    - Your SSPanel url. Make sure it start at "http" or "https"
    
- key
    - Your SSPanel's mu_key. Check from website root ./config/.config.php
    
- node_id
    - The node you want to build
    
- alert_id
    - Make sure alertId is equal to thr node configuration at your SSPanel, the wrong value will cause connection failure or memory leak
    
- in_tags
    - An array includes Tags which v2ray-core inbound you want to add users
    
- api_address
    - V2ray-core api address, normally it will be "127.0.0.1"
    
- api_port
    - V2ray-core api port, normally it will be "10085". You can change it in your v2ray-core config json.
    
- sync_interval
    - Interval time(second) in two synchronization.
    
- fail_delay
    - Retry delay time(second) if synchronization failure.