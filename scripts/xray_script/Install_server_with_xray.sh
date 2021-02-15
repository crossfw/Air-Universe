#!/bin/bash

VERSION="0.3.2"

panelConfig() {
  echo "Air-Universe 0.3.2  + Xray 1.3.0 Installation"
  echo "########Air-Universe config#######\n"
  read -p "Enter node_ids, (eg 1,2,3): " nIds
  read -p "Enter sspanel domain(https://): " pUrl
  read -p "Enter panel token: " nKey

  apt-get update
  apt-get install cron wget ca-certificates -y
}

download() {
  v2ray_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/proxy-core/xray-1_3_0"
  airuniverse_url="https://github.com/crossfw/Air-Universe/releases/download/v${VERSION}/Air-Universe-linux-amd64"
  v2ray_json_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/xray_json/multiIn.json"
  start_script_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/xray_script/Start_AU_with_xray.sh"

  wget -N --no-check-certificate ${v2ray_url} -O /usr/bin/au/xr
  wget -N --no-check-certificate ${v2ray_json_url} -O /etc/au/xr.json
  wget -N --no-check-certificate ${airuniverse_url} -O /usr/bin/au/au
  wget -N --no-check-certificate ${start_script_url} -O /usr/bin/au/run.sh
}

makeConfig() {
  cat >>/etc/au/au.json <<EOF
  {
    "panel": {
      "url": "https://${pUrl}",
      "key": "${nKey}",
      "node_ids": [${nIds}]
    },
    "proxy": {
      "log_path": "/var/log/xr.log"
    }
  }
EOF

}


keepalive() {
  cat >>/usr/bin/au/keepalive.sh <<EOF
#!/bin/bash
au=\$(ps -ef | grep "au -c" | grep -v grep | wc -l)
if [ \$au -eq 0 ];then
    bash  /usr/bin/au/run.sh restart
else
   echo "au ok"
fi

v2=\$(ps -ef | grep "xr -c" | grep -v grep | wc -l)
if [ \$v2 -eq 0 ];then
    bash  /usr/bin/au/run.sh restart
else
   echo "v2 ok"
fi
EOF

}

mkdir /usr/bin/au/
mkdir /etc/au/
panelConfig
download
makeConfig
keepalive
chmod +x /usr/bin/au/*
echo '*/1 * * * * /usr/bin/au/keepalive.sh'  >> /var/spool/cron/crontabs/root
echo '0 6 * * * cat /dev/null > /var/log/au.log'  >> /var/spool/cron/crontabs/root
chown root:crontab /var/spool/cron/crontabs/root
chmod 600 /var/spool/cron/crontabs/root
/bin/bash /usr/bin/au/keepalive.sh
