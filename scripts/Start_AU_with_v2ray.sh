#!/bin/bash

start() {
  nohup /usr/bin/au/v2 -c /etc/au/v2.json >/dev/null 2>&1 &
  sleep 5
  nohup /usr/bin/au/au -c /etc/au/au.json >/var/log/au.log 2>&1 &
}

stop() {
  eval $(ps -ef | grep "v2 -c" | awk '{print "kill "$2}')
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
