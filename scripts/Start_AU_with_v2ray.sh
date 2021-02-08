#!/bin/bash

eval $(ps -ef | grep "v2 -c" | awk '{print "kill "$2}')
eval $(ps -ef | grep "au -c" | awk '{print "kill "$2}')

nohup ./v2 -c v2.json > ./pr.log 2>&1 &
sleep 5
nohup ./au -c au.json > ./ctl.log 2>&1 &
