#!/usr/bin/env bash

if [ $# -lt 2 ];then
    echo "[usage] bash watcher.sh <container_id> "
fi

container_id=$1
echo "container id is $container_id"
container_pid=$(sudo docker inspect --format '{{.State.Pid}}' "$container_id")

echo "container pid is $container_pid"

showNamespace() {
    sudo ls -l /proc/$container_pid/ns/  | awk -F ' ' '{print $11}'
}

showNamespacePid() {
cat /proc/$container_pid/status | grep NSpid
}
showCgroup(){
  cat /proc/$container_pid/cgroup
}

showNamespace
showNamespacePid
showCgroup