# ha-api

#### 介绍

HA管理平台后端REST API服务。
另外一个相关项目为[HA-web](https://gitee.com/openeuler/ha-web)，提供HA管理平台前端UI。

## 界面截图

TODO: 待添加截图

## 特性

易于使用.

## 文档

以下文档可帮助你了解该项目：
 - [构建文档](./docs/build.md)
 - [架构文档](./docs/architecture.md)
 - [API文档](./docs/api.md)

# 安装

该HA管理平台后端运行需要安装HA软件。

openEuler 20.03 LTS SP1系统中，你可以直接通过yum安装HA软件。

```
[root@ha1~]# yum install corosync pacemaker pcs fence-agents fence-virt corosync-qdevice sbd drbd drbd-utils
```

其他操作系统中，你可能需要自行编译HA软件并安装。

你可以在[HA安装文档](./docs/ha_install.md)和[ha-api构建文档](./docs/build.md)中了解到更多信息

## 贡献

ha-api使用golang开发。我们使用[Beego框架](https://beego.me/)来构建高性能、可靠的管理平台。欢迎任何人进行贡献。如果你在使用或者开发当中有任何问题，你也可以通过提交issue的方式联系我们。


## 码云特性

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5.  码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
