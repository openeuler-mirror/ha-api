# ha-api

## Description
ha-api is the back-end REST API for the HA management platform.

Another project is [HA-web](https://gitee.com/openeuler/ha-web), which is the front-end UI of the HA management platform. 

## Screenshots

To be added.

## Features

Easy to use.

## Documents

The following documents are provided to help you understand the project:
 - [Build Document](./docs/build_en.md)
 - [Architecture Document](./docs/architecture_en.md)
 - [API Document](./docs/api_en.md)

## Installation 

The HA software is required to run the ha-api.

On openEuler 20.03 LTS SP1, you can install the HA software using Yum:

```
[root@ha1~]# yum install corosync pacemaker pcs fence-agents fence-virt corosync-qdevice sbd drbd drbd-utils
```

For other OSs, you may need to compile the HA software and then install it.
For more details, see [install documents](./docs/install_en.md) and [build documents](./docs/build_en.md).

## Contribution

ha-api is developed with Golang. We use the [Beego framework](https://beego.me/) to develop a management platform featuring high performance and solid reliability. All contributions are welcomed. If you have any problems in using or developing, please contact us by opening an issue.


## Gitee Features

1.	Use Readme_XXX.md to mark README files with different languages, such as Readme_en.md and Readme_zh.md.
2.	Gitee blog: blog.gitee.com
3.	You can visit https://gitee.com/explore to learn about excellent open source projects on Gitee.
4.	GVP is short for Gitee Most Valuable Project.
5.	User manual provided by Gitee: https://gitee.com/help
6.	Gitee Cover People: https://gitee.com/gitee-stars/
