#!/bin/bash

# 参数：@项目名称 project @路径 path
project=$1
serverName="go-hub-$1"
path=$2

# 停止服务
ps aux | grep "$serverName" | awk '{print $2}' | xargs kill -9
# 启动服务
cd "$path$project"
pwd &&
  export PATH=${GOROOT}/bin:$PATH &&
  go version &&
  go mod tidy
echo "start server"
nohup go run cmd/clientserver/main.go -n "$serverName" >run.log 2>&1 &
echo "start success"
