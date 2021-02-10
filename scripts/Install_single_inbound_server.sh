#!/bin/bash

VERSION="0.2.0"


panelConfig(){
   echo "########Air-Universe config#######\n"
    read -p "Enter node_id:" nId
    read -p "Enter sspanel domain(https://):" pUrl
    read -p "Enter panel token:" nKey
    read -p "Enter v2ray out port:" vOut

    apt-get update
    apt-get install cron wget ca-certificates -y
}


download(){
  v2ray_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/v2ray-core/v2ray-4_34_0"
  airuniverse_url="https://github.com/crossfw/Air-Universe/releases/download/v${VERSION}/Air-Universe-linux-amd64"
  v2ray_json_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/v2ray-core_json/Single.json"
#  start_script_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/Start_AU_with_v2ray.sh"

  wget -N --no-check-certificate ${v2ray_url} -O /usr/bin/au/v2
  wget -N --no-check-certificate ${v2ray_json_url} -O /etc/au/v2.json
  wget -N --no-check-certificate ${airuniverse_url} -O /usr/bin/au/au
  wget -N --no-check-certificate ${start_script_url} -O /usr/bin/au/run.sh

  chmod +x /usr/bin/au/au
  chmod +x /usr/bin/au/v2
#  chmod +x /usr/bin/au/run.sh
}

makeConfig(){
cat >> /etc/au/au.json << EOF
{
"panel": {
  "url": "https://${pUrl}",
  "key": "${nKey}",
  "node_ids": [${nId}]
},
"proxy": {
  "log_path": "/var/log/au-v.log"
}
}
EOF

sed  -i "s/11071/${vOut}/g"  /etc/au/v2.json

}

createV2Service(){
cat >> /lib/systemd/system/v2.service << EOF
[Unit]
Description=Air-Universe Service
After=network.target

[Service]
Type=simple
User=nobody
Restart=on-failure
RestartSec=5s
ExecStart=/usr/bin/au/v2 -c /etc/au/v2.json

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload;
systemctl enable v2;
systemctl start v2;
}

createAUService(){
cat >> /lib/systemd/system/au.service << EOF
[Unit]
Description=Air-Universe Service
After=network.target

[Service]
Type=simple
User=nobody
Restart=on-failure
RestartSec=5s
ExecStart=/usr/bin/au/au -c /etc/au/au.json

[Install]
WantedBy=multi-user.target
EOF

chmod -R 777 /var/log/
systemctl daemon-reload;
systemctl enable au;
systemctl start au;
}


mkdir /usr/bin/au/
mkdir /etc/au/
panelConfig;
download;
makeConfig;
createV2Service;
createAUService;
