# Author: d1y<chenhonzhou@gmail.com>
# 编写时间: 2020-10-31

# 什么是 `yy` ??
# 用的是这个脚本: https://github.com/machinebox/appify
# `windows` 同理

# git clone https://github.com/machinebox/appify
# cd appify
# go build -o ~/go/bin/yy .

cd build

rm -rf *
rm -rf /Applications/inet.app

go build -o inet ..

yy -name "inet" -icon ../icon/icon_1024.png inet

cp -rf inet.app /Applications/inet.app