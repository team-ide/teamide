#!/bin/bash

# 配置要启动关闭的脚本名
process_name="teamide-node"

# 添加启动命令
function start(){
    echo "start..."

    chmod +x $process_name
    nohup ./$process_name -id teamide-node -address :51091 -token teamide-token > log/start.log 2>&1 &

    echo "start successful"
    return 0
}

# 添加停止命令
function stop(){
    echo "stop..."

    ps aux |grep $process_name |grep -v grep |awk '{print "kill -9 " $2}'|sh

    echo "stop successful"
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
*)
    echo "请输入: start, stop, restart"
    ;;
esac
