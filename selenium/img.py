# !/usr/bin/env python3
# -*- coding:utf-8 -*-
#

from __future__ import division
import time
import sys
from selenium import webdriver


driver = webdriver.Remote(
    command_executor='http://127.0.0.1:4444',
    options=webdriver.FirefoxOptions()
)
driver.delete_all_cookies()

driver.set_window_size(1366, 768)


user = sys.argv[1]
pw = sys.argv[2]

sort = ["live", "cdn"]
dm = ["day=14"]

for s in iter(sort):
    other = ""
    if s == "cdn":
        other = "cdnOtherGB=100"
    for d in iter(dm):
        time.sleep(1)
        driver.get(
            f"http://{user}:{pw}@127.0.0.1:8174/{s}?downloadImg=true&{other}&{d}"
        )
        time.sleep(2)

driver.quit()
