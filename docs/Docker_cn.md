# Docker Install

此方法暂时废弃
## 使用方法
```
sudo docker run --network host -v /PATH2CONFIG:/usr/local/etc/au crossfw/airu
```
其中在 /PATH2CONFIG 下创建au.json文件(Air-Universe)配置文件

Example:
```
sudo docker run --network host -v /root/aucfg:/usr/local/etc/au crossfw/airu
```