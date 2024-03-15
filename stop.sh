#!/bin/bash

# 参数：@项目名称 project
pwd
# 停止服务
serverName="go-hub-$1"
ps aux | grep "$serverName" | awk '{print $2}' | xargs kill -9
