#!/bin/sh

# 移至脚本目录
cd `dirname $0`

echo `pwd`

export LD_LIBRARY_PATH=$(pwd)/lib:$LD_LIBRARY_PATH
echo $LD_LIBRARY_PATH

mkdir -p log

# 添加启动命令
function start(){
    echo "start..."

    chmod +x teamide
    echo " ldd teamide"
    ldd teamide
    nohup ./teamide TeamIDE-Server > log/start.log 2>&1 &

    echo "start successful"
    return 0
}

# 添加停止命令
function stop(){
    echo "stop..."

    ps aux |grep TeamIDE-Server |grep -v grep |awk '{print "kill -9 " $2}'|sh

    echo "stop successful"
    return 0
}

function version(){
    echo "version..."
    chmod +x teamide
    ./teamide -v
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
    echo "请输入: start, stop, restart"
    ;;
esac