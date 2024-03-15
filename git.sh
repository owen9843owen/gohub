#!/bin/bash

# 思路 创建目录，然后clone并切换分支
# 参数 @仓库地址 @分支 @本地路径 @项目文件夹
pwd
# 获取最新版本信息
# 配置go版本
repository=$1
branch=$2
path=$3
project=$4
# 检查运行目录，没有则创建
if [ ! -d "$path" ]; then
  mkdir -p "$path"
fi
cd "$path"
# 检查仓库，没有则clone代码
if [ ! -d "$path$project" ]; then
  git clone "$repository"
fi
# 进入项目，切换分支，更新代码
cd "$path$project"
git checkout "$branch"
git pull --all
git rev-parse HEAD
