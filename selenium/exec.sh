#!/usr/bin/env bash
# 0 2 * * 1 bash /pd/dev/py/auto-check/exec.sh


# 执行
source /etc/profile
docker restart selenium
sleep 5

# 文件目录绝对路径
dir_path=$(dirname $0)
if [[ $dir_path == "." ]] ;then
    dir_path=$(pwd)
fi

echo "" > ${dir_path}/log
python3 ${dir_path}/selenium.py $1 $2 > ${dir_path}/log 2>&1 \
&& docker stop selenium
