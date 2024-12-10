#!/bin/bash

# 检查输入参数
if [ $# -lt 2 ]; then
    echo "Usage: $0 <pid> <namespace_type>"
    echo "Example: $0 12345 pid"
    exit 1
fi

TARGET_PID=$1
NS_TYPE=$2

# 获取目标 Namespace ID
NS_ID=$(sudo readlink "/proc/$TARGET_PID/ns/$NS_TYPE" |  sed -n 's/.*\[\([0-9]*\)\].*/\1/p' )
#echo "get ns_id $NS_ID"
if [ -z "$NS_ID" ]; then
    echo "Error: Failed to retrieve Namespace ID for PID $TARGET_PID and type $NS_TYPE."
    exit 1
fi

echo "Namespace Type: $NS_TYPE"
#echo "Namespace ID: $NS_ID"
echo "peers in the same namespace:"
for pid in $(ls /proc | grep '^[0-9]\+$'); do
        CURRENT_NS_ID=$(sudo readlink /proc/$pid/ns/$NS_TYPE | grep -o '[0-9]\+')
        if [ "$CURRENT_NS_ID" == "$NS_ID" ]; then
            PROCESS_INFO=$(ps -p $pid -o pid,ppid,user,comm --no-headers 2>/dev/null)
            echo "$PROCESS_INFO"
        fi
done