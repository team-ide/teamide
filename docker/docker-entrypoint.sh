#!/bin/sh

# 移至脚本目录
cd `dirname $0`

echo `pwd`

export LD_LIBRARY_PATH=$(pwd)/lib:$LD_LIBRARY_PATH
echo $LD_LIBRARY_PATH
chmod +x teamide

./teamide