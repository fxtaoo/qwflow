# 说明

每周一发送图片邮件报表，外部依赖 selenium 访问 selenium/standalone-firefox 生成图片图片文件

程序再将本地图片编码 base64 以邮件发送

具体实现：
1. 系统 crontab 定时执行脚本 exec.sh，参数 $1 $2 分别为账号密码，程序是凌晨 3 点，crontab 需要配置为 3 点前
2. exec.sh 执行同文件夹下 selenium.py 控制 selenium/standalone-firefox 浏览器访问
3. 网址包含 downloadImg=true 会自动下载保存图片至 selenium/standalone-firefox 容器下 /home/seluser/Downloads

img 文件夹需要挂载到 fxtaoo/qwflow selenium/standalone-firefox 两个容器，注意路径与读写权限

需要新安装包 selenium
pip3 install -i https://mirrors.cloud.tencent.com/pypi/simple selenium==4.4.3
