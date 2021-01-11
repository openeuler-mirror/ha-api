# ha-api

#### 介绍
API interface for HA management

#### 软件架构
软件架构说明


#### 安装教程

环境需求：至少两台已经安装好openEuler 20.03 LTS-SP1的物理机/虚拟机

***注：以下操作两台主机均需要操作，现以ha1主机为例***
1.  修改主机名称及 /etc/hosts文件
在使用HA软件之前，需要确认修改主机名并将所有主机名写入/etc/hosts中。

```
[root@ha1 ~]# hostnamectl set-hostname ha1
[root@ha1 ~]# vim /etc/hosts
10.1.80.20 ha1
10.1.80.21 ha2
```

2.  关闭防火墙
```
[root@ha1 ~]# systemctl stop firewalld
```
修改SELINUX状态为disabled
```
[root@ha1 ~]# vim /etc/selinux/config
SELINUX=disabled
SELINUXTYPE=targeted
```

3.  安装软件包
成功安装系统后，会默认配置好yum源，文件位置存放在/etc/yum.repos.d/openEuler.repo文件中，软件包会用到以下源:
```
[OS]
name=OS
baseurl=http://repo.openeuler.org/openEuler-20.03-LTS-SP1/OS/$basearch/
enabled=1
gpgcheck=1
gpgkey=http://repo.openeuler.org/openEuler-20.03-LTS-SP1/OS/$basearch/RPM-GPG-KEY-openEuler

[everything]
name=everything
baseurl=http://repo.openeuler.org/openEuler-20.03-LTS-SP1/everything/$basearch/
enabled=1
gpgcheck=1
gpgkey=http://repo.openeuler.org/openEuler-20.03-LTS-SP1/everything/$basearch/RPM-GPG-KEY-openEuler

[EPOL]
name=EPOL
baseurl=http://repo.openeuler.org/openEuler-20.03-LTS-SP1/EPOL/$basearch/
enabled=1
gpgcheck=1
gpgkey=http://repo.openeuler.org/openEuler-20.03-LTS-SP1/OS/$basearch/RPM-GPG-KEY-openEuler
```
使用以下命令安装HA软件包
```
[root@ha1~]# yum install corosync pacemaker pcs fence-agents fence-virt corosync-qdevice sbd drbd drbd-utils
```
4. 设置hacluster用户密码
```
[root@ha1~]# passwd hacluster
```
5. 修改corosync.conf文件
```
totem {
        version: 2
        cluster_name: hacluster
         crypto_cipher: none
        crypto_hash: none
}
logging {         
        fileline: off
        to_stderr: yes
        to_logfile: yes
        logfile: /var/log/cluster/corosync.log
        to_syslog: yes
        debug: on
       logger_subsys {
               subsys: QUORUM
               debug: on
        }
}
quorum {
           provider: corosync_votequorum
           expected_votes: 2
           two_node: 1
       }
nodelist {
       node {
               name: ha1
               nodeid: 1
               ring0_addr: 10.1.80.21
               }
        node {
               name: ha2
               nodeid: 2
               ring0_addr: 10.1.80.22
               }
        }
```
6. 启动服务
启动以下服务：
```
[root@ha1~]# systemctl start pcsd
[root@ha1~]# systemctl start pacemaker
[root@ha1~]# systemctl start corosync
```
7. 节点鉴权
***注：一个节点上执行即可***
```
[root@ha1~]# pcs host auth ha1 ha2
```
8. 访问前端管理平台
在浏览器中直接访问`https://IP:2224`即可。用户名为`hacluster`，密码为该用户在主机上设置的密码。
#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 码云特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5.  码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
