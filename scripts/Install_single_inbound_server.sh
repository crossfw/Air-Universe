#!/bin/bash


download(){
  v2ray_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/v2ray-core/v2ray-4_34_0'
  airuniverse_url='https://github.com/crossfw/Air-Universe/releases/download/v0.1.0/Air-Universe-linux-amd64'
  v2ray_json_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/v2ray-core_json/Single.json'
  airuniverse_json_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/Air-Universe_json/Air-Universe_full.json'
  start_script_url='https://raw.githubusercontent.com/crossfw/Air-Universe/master/scripts/Start_AU_with_v2ray.sh'
  wget -N --no-check-certificate ${v2ray_url} -O ./v2
  wget -N --no-check-certificate ${v2ray_json_url} -O ./v2.json
  wget -N --no-check-certificate ${airuniverse_url} -O ./au
  wget -N --no-check-certificate ${airuniverse_json_url} -O ./au.json
  wget -N --no-check-certificate ${start_script_url} -O ./run.sh
  chmod +x ./v2
  chmod +x ./au
  chmod +x ./run.sh
}

mkdir ./au
cd au
download;
