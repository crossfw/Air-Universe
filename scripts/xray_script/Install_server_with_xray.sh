#!/bin/bash

VERSION=""
APP_PATH="/usr/local/bin/"
CONFIG_PATH="/usr/local/etc/au/"

create_folders() {
  if [[ ! -e "${APP_PATH}" ]]; then
    mkdir "${APP_PATH}"
  fi
  if [[ ! -e "${CONFIG_PATH}" ]]; then
    mkdir "${CONFIG_PATH}"
  fi

}

panelConfig() {
  echo "Air-Universe $VERSION + Xray"
  echo "########Air-Universe config#######"
  read -r -p "Enter node_ids, (eg 1,2,3): " nIds
  read -r -p "Enter panel domain(Include https:// or http://): " pUrl
  read -r -p "Enter panel token: " nKey
  echo && echo -e "Choose panel type:
  1. SSPanel
  2. V2board"
  read -r -p "Choose panel type: " panelnum
  if [ "$panelnum" == "1" ]; then
    panelType="sspanel"
  fi
  if [ "$panelnum" == "2" ]; then
    panelType="v2board"
    read -r -p "Enter nodes type, (eg \"vmess\",\"ss\")(DON'T FORGET '\"'): " nType
    read -r -p "Enter nodes enable receive proxy protocol, (eg true, false) enter means all false: " npp
  fi
}

check_root() {
  [[ $EUID != 0 ]] && echo -e "${Error} 当前非ROOT账号(或没有ROOT权限)，无法继续操作，请更换ROOT账号或使用 ${Green_background_prefix}sudo su${Font_color_suffix} 命令获取临时ROOT权限（执行后可能会提示输入当前账号的密码）。" && exit 1
}
check_sys() {
  if [[ -f /etc/redhat-release ]]; then
    release="centos"
  elif cat /etc/issue | grep -q -E -i "debian"; then
    release="debian"
  elif cat /etc/issue | grep -q -E -i "ubuntu"; then
    release="ubuntu"
  elif cat /etc/issue | grep -q -E -i "centos|red hat|redhat"; then
    release="centos"
  elif cat /proc/version | grep -q -E -i "debian"; then
    release="debian"
  elif cat /proc/version | grep -q -E -i "ubuntu"; then
    release="ubuntu"
  elif cat /proc/version | grep -q -E -i "centos|red hat|redhat"; then
    release="centos"
  fi
  bit=$(uname -m)
}
Installation_dependency() {
  if [[ ${release} == "centos" ]]; then
    yum update
    yum install -y gzip ca-certificates curl
  else
    apt-get update
    apt-get install -y ca-certificates curl
  fi
  cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
  mkdir /var/log/au
  chown -R nobody /var/log/au
}
download() {
  mkdir /usr/local/etc/au/
  airuniverse_url="https://github.com/crossfw/Air-Universe/releases/download/${VERSION}/Air-Universe-linux-amd64"
  xray_json_url="https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/xray_json/multiIn.json"

  mv /usr/local/etc/xray/config.json /usr/local/etc/xray/config.json.bak
  wget -N --no-check-certificate ${xray_json_url} -O /usr/local/etc/xray/config.json
  wget -N --no-check-certificate ${airuniverse_url} -O /usr/local/bin/au

  chmod +x /usr/local/bin/au


}
get_latest_version() {
  # Get Xray latest release version number
  local tmp_file
  tmp_file="$(mktemp)"
  if ! curl -x "${PROXY}" -sS -H "Accept: application/vnd.github.v3+json" -o "$tmp_file" 'https://api.github.com/repos/crossfw/Air-Universe/releases/latest'; then
    "rm" "$tmp_file"
    echo 'error: Failed to get release list, please check your network.'
    exit 1
  fi
  RELEASE_LATEST="$(sed 'y/,/\n/' "$tmp_file" | grep 'tag_name' | awk -F '"' '{print $4}')"
  if [[ -z "$RELEASE_LATEST" ]]; then
    if grep -q "API rate limit exceeded" "$tmp_file"; then
      echo "error: github API rate limit exceeded"
    else
      echo "error: Failed to get the latest release version."
      echo "Welcome bug report:https://github.com/crossfw/Air-Universe/issues"
    fi
    "rm" "$tmp_file"
    exit 1
  fi
  "rm" "$tmp_file"
  VERSION="v${RELEASE_LATEST#v}"
}
makeConfig() {
  mkdir -p /usr/lib/systemd/system/
  cat >/usr/local/etc/au/au.json <<EOF
{
  "panel": {
    "type": "${panelType}",
    "url": "${pUrl}",
    "key": "${nKey}",
    "node_ids": [${nIds}],
    "nodes_type": [${nType}],
    "nodes_proxy_protocol": [${npp}]
  },
  "proxy": {
    "type":"xray",
    "log_path": "/var/log/au/xr.log"
  }
}
EOF
}

createService() {
  cat >/usr/lib/systemd/system/au.service <<EOF
[Unit]
Description=Air-Universe - main Service
After=network.target
Wants=xray.service

[Service]
Type=simple
User=nobody
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/au -c /usr/local/etc/au/au.json > /var/log/au/au.log

[Install]
WantedBy=multi-user.target
EOF
}
install_xray(){
  bash -c "$(curl -L https://github.com/crossfw/Xray-install/raw/main/install-release.sh)" @ install
}

check_root
check_sys
Installation_dependency
install_xray
get_latest_version
panelConfig
download
makeConfig
createService

systemctl enable au
systemctl restart xray
systemctl start au
