#!/bin/sh

# 移至脚本目录
cd `dirname $0`

echo `pwd`

cd ../

templateDir=./electron-template
# 安装 electron-template

git clone https://github.com/team-ide/electron-template
rm -rf $templateDir/.git

# 复制 应用信息 package.json 到 release/app
echo 'cp package.json'
cp -rf package.json $templateDir/release/app/package.json

echo 'app package.json info'
cat $templateDir/release/app/package.json

# 复制 应用配置 config.ts 到 src/main
echo 'cp config.ts'
cp -rf desktop/config.ts $templateDir/src/main/config.ts

echo 'config.ts info'
cat $templateDir/src/main/config.ts

# 设置 应用 变量
productName='TeamIDE'
publisherName='ZhuLiang'
publishProvider='github'
publishOwner='team-ide'
publishRepo='teamide'

echo 'set productName='$productName
echo 'set publisherName='$publisherName
echo 'set publishProvider='$publishProvider
echo 'set publishOwner='$publishOwner
echo 'set publishRepo='$publishRepo

# 设置包相关信息

# 设置 项目名称
echo 'replace productName'
sed -i 's/<productName>/'$productName'/g' $templateDir/package.json

# 设置 项目 发布者
echo 'replace publisherName'
sed -i 's/<publisherName>/'$publisherName'/g' $templateDir/package.json

# 设置 项目 发布信息
echo 'replace publish'
sed -i 's/"<publish>"/{"provider": "'$publishProvider'","owner": "'$publishOwner'","repo": "'$publishRepo'"}/g' $templateDir/package.json

echo 'package.json info'
cat $templateDir/package.json