# 架构

该文件描述ha-api管理平台项目后端架构。

## 项目架构

TODO: 待补充图片

总体来讲，ha-api后端封装了很多的HA软件管理命令行，例如pcs, crm_xxx, cibamdin等。然后通过ha-web项目提供一个易于使用的前端界面来监控和管理HA集群。


## 代码结构

```
--
 |- controllers
 |- models
 |- routers
 |- services
 |- settings
 |- utils
 |- views/static

```

controllers: Beego框架中的REST api处理controller。
models: 封装HA集群管理命令。
routers: 绑定URL和controller。
services: 当前只有session服务.
settings: 应用配置.
utils: 公共工具类.
views: 静态文件.

