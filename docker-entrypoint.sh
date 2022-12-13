#!/bin/sh


set -ex;

# 移至脚本目录
cd `dirname $0`

echo `pwd`

export LD_LIBRARY_PATH=/opt/teamide/lib/:/lib:/lib64:/usr/lib:/usr/lib64:$LD_LIBRARY_PATH
echo $LD_LIBRARY_PATH
chmod +x teamide

./teamide