#!/bin/bash

panelConfig(){
   echo "########Air-Universe config#######\n"
    read -p "Enter node_id:" nId
    read -p "Enter sspanel domain(https://):" pUrl
    read -p "Enter panel token:" nKey

    apt-get update
    apt-get install cron wget -y
}


download(){
  v2ray_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/v2ray-core/v2ray-4_34_0'
  airuniverse_url='https://github.com/crossfw/Air-Universe/releases/download/v0.1.1/Air-Universe-linux-amd64'
  v2ray_json_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/v2ray-core_json/Single.json'
  start_script_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/Start_AU_with_v2ray.sh'

  wget -N --no-check-certificate ${v2ray_url} -O /usr/bin/au/v2
  wget -N --no-check-certificate ${v2ray_json_url} -O /etc/au/v2.json
  wget -N --no-check-certificate ${airuniverse_url} -O /usr/bin/au/au
  wget -N --no-check-certificate ${start_script_url} -O /usr/bin/au/run.sh

  chmod +x /usr/bin/au/au
  chmod +x /usr/bin/au/v2
  chmod +x /usr/bin/au/run.sh
}

makeAUConfig(){
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
}

createService(){
  cat >> /lib/systemd/system/au.service << EOF
  [Unit]
Description=Air-Universe Service
After=network.target

[Service]
Type=simple
User=nobody
Restart=on-failure
RestartSec=5s
ExecStart=/usr/bin/au/run.sh

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload;
systemctl enable au;
systemctl start au;
}


mkdir /usr/bin/au/
mkdir /etc/au/
panelConfig;
makeAUConfig;
download;
createService;

