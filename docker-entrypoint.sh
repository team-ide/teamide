#!/bin/sh


set -ex;

# 移至脚本目录
cd `dirname $0`

echo `pwd`

export LD_LIBRARY_PATH=/opt/teamide/lib/
echo $LD_LIBRARY_PATH
chmod +x teamide

./teamide