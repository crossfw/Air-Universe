#!/bin/bash

start() {
  nohup /usr/bin/au/xr -c /etc/au/xr.json >/dev/null 2>&1 &
  sleep 5
  nohup /usr/bin/au/au -c /etc/au/au.json >/var/log/au.log 2>&1 &
}

stop() {
  eval $(ps -ef | grep "xr -c" | awk '{print "kill "$2}')
  eval $(ps -ef | grep "au -c" | awk '{print "kill "$2}')
}

restart() {
  stop
  start
}
case "$1" in
start)
  start
  ;;
stop)
  stop
  ;;
restart)
  stop
  start
  ;;
*)
  echo "start|stop|restart"
  ;;
esac
