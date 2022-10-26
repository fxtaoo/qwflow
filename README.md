# qwflow
使用 echarts 展示七牛网宿直播 cdn 带宽流量

## 说明

使用 gin template + echarts 展网宿七牛带宽折线图，流量饼图。

只需要配置好数据库、七牛网宿密钥，自动收集展示七牛网宿在用的直播、cdn 相关带宽折线图，流量饼图。

每日定时从七牛网宿获取直播 cdn 相关带宽，流量原始数据，处理后存入数据库中，访问 web 时，从数据库获取指定范围日期数据，处理转换为带宽折线图、流量饼图。

`conf-example.json` 为配置文件示例，配置数据库、七牛网宿密钥，使用需要更名为 `conf.json`。

`flow.sql` 为 mysql 表结构，数据的存储与读取依赖固定的表结构。

## 待办列表

- [x] 七牛自动动获取在用直播空间、cdn 域名列表，无需手动配置维护更新
- [x] 七牛直播空间 cdn 相关带宽、流量，数据获取、处理、存储、图表展示
- [ ] 七牛多个直播空间一天同一个时刻的带宽之和峰值
- [ ] 七牛多个 cdn 一天中同一个时刻的带宽之和峰值
- [ ] 七牛全局加速 cdn 带宽、流量，数据获取、处理、存储、图表展示
- [x] 网宿直播域名、cdn 相关带宽、流量，数据获取、处理、存储、图表展示
- [x] 网宿属于同一业务的多个直播域名，一天同一个时刻带宽之和峰值
- [ ] 网宿多个业务一天同一个时刻带宽之和峰值
- [x] 每日定时从七牛网宿获取直播 cdn 相关带宽，流量原始数据，处理后存入数据库中
- [x] 七牛网宿汇总展示
- [x] 选择日期展示流量
- [ ] 登陆权限验证
- [ ] 容器方式运行
