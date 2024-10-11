#!/bin/bash

# 移至脚本目录
cd `dirname $0`

echo `pwd`

mkdir -p log

# 添加启动命令
function start(){
    echo "start..."

    chmod +x node
    echo " ldd node"
    ldd node
    nohup ./node -address=0.0.0.0:2013 -id=TeamIDE-Node-HW-Cloud -token=x1x2x3x4 > log/start.log 2>&1 &

    echo "start successful"
    return 0
}

# 添加停止命令
function stop(){
    echo "stop..."

    ps aux |grep TeamIDE-Node |grep -v grep |awk '{print "kill -9 " $2}'|sh

    echo "stop successful"
    return 0
}

function version(){
    echo "version..."
    chmod +x node
    ./node -v
    return 0
}

case $1 in
"start")
    start
    ;;
"stop")
    stop
    ;;
"restart")
    stop && start
    ;;
"v")
    version
    ;;
"version")
    version
    ;;
*)
    echo "请输入: start, stop, restart, version"
    ;;
esac
