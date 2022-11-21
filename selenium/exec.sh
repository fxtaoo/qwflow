#!/usr/bin/env bash

# 执行
source /etc/profile

# 文件目录绝对路径
dir_path=$(dirname $0)
if [[ $dir_path == "." ]] ;then
    dir_path=$(pwd)
fi

if [ "$1" == '' ] || [ "$2" == '' ] ;then
    echo "缺少 BasicAuth 账号密码，分别为 \$1 \$2"
    exit 1
fi

docker restart selenium
sleep 5

echo "" > ${dir_path}/log

rm -rf ${dir_path}/img/*

python3 ${dir_path}/img.py $1 $2 > ${dir_path}/log 2>&1

docker stop selenium
