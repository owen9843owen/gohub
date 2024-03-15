#!/bin/bash
pwd
# 获取最新版本信息
# 配置go版本
serverName="go-hub-$1"
commonRepository=$2
path=$3
repository=$4
branch=$5

# 运行目录
echo "$serverName"
echo "$commonRepository"
echo "$path"
echo "$repository"
echo "$branch"

mkdir -p "$path"
cd "$path"
# stop
ps aux | grep "$serverName" | awk '{print $2}' | xargs kill -9
# clone代码，切换分支，更新代码
git clone "$commonRepository"
git clone "$repository"
cd "$1" &&
git checkout "$branch" &&
git pull --all &&
git rev-parse HEAD
# 启动服务
pwd &&
export PATH=${GOROOT}/bin:$PATH &&
go version
pwd &&
echo "version" &&
go mod tidy &&
echo "success" &&
nohup go run cmd/clientserver/main.go -n "$serverName" >run.log 2>&1 &
