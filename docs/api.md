# API接口文档

## 1 API说明

### 1.1 API概述

该文档定义HA管理平台后端接口。

### 1.2 状态码

| 状态码 | 说明                           |
| :----: | :----------------------------- |
|  2xx   | 请求正常处理并返回             |
|  3xx   | 重定向，请求的资源位置发生变化 |
|  4xx   | 客户端发送请求有误             |
|  5xx   | 服务端错误                     |

### 1.3 错误返回

所有API均采用统一的错误返回格式，具体描述如下：

响应格式：json
响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| error                         | string                             | 返回的错误信息                        |

错误响应示例：
```
{
    "action":false,
    "error":"error info"
}
```


### 1.4 修改记录

| 序号 | 更新内容 |    日期    |
| :--: | :------- | :--------: |
|  1   | 初稿     | 2021.01.11 |


<div STYLE="page-break-after: always;"></div>



## 2 API接口定义

### 2.1 集群

#### 2.1.1 获取集群属性

URI：/api/v1/haclusters/1

Method：GET

请求参数：

无


响应格式：json
响应参数

| 参数名称                          | 参数类型                               | 参数说明                             |
| :------------------------------ - | :------------------------------------ | ------------------------------------ |
| action                            | bool                                  | 产品总条数                           |
| data                              | object                                | 返回的正常结果                       |
| error                             | string                                | 返回的错误信息                       |
| parameters                        | object                                | 各个组件的参数详情                   |
| shortdesc                         | string                                | 简短描述                             |
| version                           | string                                | 版本信息                             |
| nodecount                         | int                                   | 挂载的节点数量                       |
| isconfig                          | bool                                  | 是否是配置表述                       |
| longdesc                          | string                                | 完整描述                             |



响应示例

```
{
    "action":true,
    "data":{
        "name":"Policy Engine",
        "parameters":{
            "stonith-timeout":{
                "name":"stonith-timeout",
                "enabled":0,
                "value":"60s",
                "content":{
                    "default":"60s",
                    "type":"time"
                },
                "shortdesc":"How long to wait for the STONITH action (reboot,on,off) to complete",
                "unique":"0",
                "longdesc":"How long to wait for the STONITH action (reboot,on,off) to complete"
            },
            "node-health-green":{
                "name":"node-health-green",
                "enabled":1,
                "value":"0",
                "content":{
                    "default":"0",
                    "type":"integer"
                },
                "shortdesc":"The score 'green' translates to in rsc_location constraints",
                "unique":"0",
                "longdesc":"Only used when node-health-strategy is set to custom or progressive."
            },
            "placement-strategy":{
                "name":"placement-strategy",
                "enabled":0,
                "value":"default",
                "content":{
                    "default":"default",
                    "values":[
                        "default",
                        "utilization",
                        "minimal",
                        "balanced"
                    ],
                    "type":"enum"
                },
                "shortdesc":"The strategy to determine resource placement",
                "unique":"0",
                "longdesc":"The strategy to determine resource placement  Allowed values: default, utilization, minimal, balanced"
            },
            "symmetric-cluster":{
                "name":"symmetric-cluster",
                "enabled":1,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"All resources can run anywhere by default",
                "unique":"0",
                "longdesc":"All resources can run anywhere by default"
            },
            "pe-input-series-max":{
                "name":"pe-input-series-max",
                "enabled":0,
                "value":"4000",
                "content":{
                    "default":"4000",
                    "type":"integer"
                },
                "shortdesc":"The number of other PE inputs to save",
                "unique":"0",
                "longdesc":"Zero to disable, -1 to store unlimited."
            },
            "maintenance-mode":{
                "name":"maintenance-mode",
                "enabled":1,
                "value":"false",
                "content":{
                    "default":"false",
                    "type":"boolean"
                },
                "shortdesc":"Should the cluster monitor resources and start/stop them as required",
                "unique":"0",
                "longdesc":"Should the cluster monitor resources and start/stop them as required"
            },
            "default-action-timeout":{
                "name":"default-action-timeout",
                "enabled":0,
                "value":"20s",
                "content":{
                    "default":"20s",
                    "type":"time"
                },
                "shortdesc":"How long to wait for actions to complete",
                "unique":"0",
                "longdesc":"How long to wait for actions to complete"
            },
            "startup-fencing":{
                "name":"startup-fencing",
                "enabled":0,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"STONITH unseen nodes",
                "unique":"0",
                "longdesc":"Advanced Use Only!  Not using the default is very unsafe!"
            },
            "node-health-yellow":{
                "name":"node-health-yellow",
                "enabled":1,
                "value":"0",
                "content":{
                    "default":"0",
                    "type":"integer"
                },
                "shortdesc":"The score 'yellow' translates to in rsc_location constraints",
                "unique":"0",
                "longdesc":"Only used when node-health-strategy is set to custom or progressive."
            },
            "start-failure-is-fatal":{
                "name":"start-failure-is-fatal",
                "enabled":1,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"Always treat start failures as fatal",
                "unique":"0",
                "longdesc":"This was the old default.  However when set to FALSE, the cluster will instead use the resource's failcount and value for resource-failure-stickiness"
            },
            "enable-startup-probes":{
                "name":"enable-startup-probes",
                "enabled":0,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"Should the cluster check for active resources during startup",
                "unique":"0",
                "longdesc":"Should the cluster check for active resources during startup"
            },
            "stop-orphan-actions":{
                "name":"stop-orphan-actions",
                "enabled":0,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"Should deleted actions be cancelled",
                "unique":"0",
                "longdesc":"Should deleted actions be cancelled"
            },
            "stop-all-resources":{
                "name":"stop-all-resources",
                "enabled":0,
                "value":"false",
                "content":{
                    "default":"false",
                    "type":"boolean"
                },
                "shortdesc":"Should the cluster stop all active resources (except those needed for fencing)",
                "unique":"0",
                "longdesc":"Should the cluster stop all active resources (except those needed for fencing)"
            },
            "default-resource-stickiness":{
                "name":"default-resource-stickiness",
                "enabled":0,
                "value":"0",
                "content":{
                    "default":"0",
                    "type":"integer"
                },
                "shortdesc":"",
                "unique":"0",
                "longdesc":""
            },
            "no-quorum-policy":{
                "name":"no-quorum-policy",
                "enabled":1,
                "value":"ignore",
                "content":{
                    "default":"stop",
                    "values":[
                        "stop",
                        "freeze",
                        "ignore",
                        "suicide"
                    ],
                    "type":"enum"
                },
                "shortdesc":"What to do when the cluster does not have quorum",
                "unique":"0",
                "longdesc":"What to do when the cluster does not have quorum  Allowed values: stop, freeze, ignore, suicide"
            },
            "node-health-red":{
                "name":"node-health-red",
                "enabled":1,
                "value":"-INFINITY",
                "content":{
                    "default":"-INFINITY",
                    "type":"integer"
                },
                "shortdesc":"The score 'red' translates to in rsc_location constraints",
                "unique":"0",
                "longdesc":"Only used when node-health-strategy is set to custom or progressive."
            },
            "batch-limit":{
                "name":"batch-limit",
                "enabled":0,
                "value":"0",
                "content":{
                    "default":"0",
                    "type":"integer"
                },
                "shortdesc":"The number of jobs that the TE is allowed to execute in parallel",
                "unique":"0",
                "longdesc":"The "correct" value will depend on the speed and load of your network and cluster nodes."
            },
            "concurrent-fencing":{
                "name":"concurrent-fencing",
                "enabled":0,
                "value":"false",
                "content":{
                    "default":"false",
                    "type":"boolean"
                },
                "shortdesc":"Allow performing fencing operations in parallel",
                "unique":"0",
                "longdesc":"Allow performing fencing operations in parallel"
            },
            "stonith-enabled":{
                "name":"stonith-enabled",
                "enabled":1,
                "value":"false",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"Failed nodes are STONITH'd",
                "unique":"0",
                "longdesc":"Failed nodes are STONITH'd"
            },
            "have-watchdog":{
                "name":"have-watchdog",
                "enabled":0,
                "value":"false",
                "content":{
                    "default":"false",
                    "type":"boolean"
                },
                "shortdesc":"Enable watchdog integration",
                "unique":"0",
                "longdesc":"Set automatically by the cluster if SBD is detected.  User configured values are ignored."
            },
            "stop-orphan-resources":{
                "name":"stop-orphan-resources",
                "enabled":0,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"Should deleted resources be stopped",
                "unique":"0",
                "longdesc":"Should deleted resources be stopped"
            },
            "stonith-action":{
                "name":"stonith-action",
                "enabled":0,
                "value":"reboot",
                "content":{
                    "default":"reboot",
                    "values":[
                        "reboot",
                        "poweroff",
                        "off"
                    ],
                    "type":"enum"
                },
                "shortdesc":"Action to send to STONITH device",
                "unique":"0",
                "longdesc":"Action to send to STONITH device  Allowed values: reboot, poweroff, off"
            },
            "pe-warn-series-max":{
                "name":"pe-warn-series-max",
                "enabled":0,
                "value":"5000",
                "content":{
                    "default":"5000",
                    "type":"integer"
                },
                "shortdesc":"The number of PE inputs resulting in WARNINGs to save",
                "unique":"0",
                "longdesc":"Zero to disable, -1 to store unlimited."
            },
            "pe-error-series-max":{
                "name":"pe-error-series-max",
                "enabled":0,
                "value":"-1",
                "content":{
                    "default":"-1",
                    "type":"integer"
                },
                "shortdesc":"The number of PE inputs resulting in ERRORs to save",
                "unique":"0",
                "longdesc":"Zero to disable, -1 to store unlimited."
            },
            "migration-limit":{
                "name":"migration-limit",
                "enabled":0,
                "value":"-1",
                "content":{
                    "default":"-1",
                    "type":"integer"
                },
                "shortdesc":"The number of migration jobs that the TE is allowed to execute in parallel on a node",
                "unique":"0",
                "longdesc":"The number of migration jobs that the TE is allowed to execute in parallel on a node"
            },
            "is-managed-default":{
                "name":"is-managed-default",
                "enabled":0,
                "value":"true",
                "content":{
                    "default":"true",
                    "type":"boolean"
                },
                "shortdesc":"Should the cluster start/stop resources as required",
                "unique":"0",
                "longdesc":"Should the cluster start/stop resources as required"
            },
            "node-health-strategy":{
                "name":"node-health-strategy",
                "enabled":1,
                "value":"none",
                "content":{
                    "default":"none",
                    "values":[
                        "none",
                        "migrate-on-red",
                        "only-green",
                        "progressive",
                        "custom"
                    ],
                    "type":"enum"
                },
                "shortdesc":"The strategy combining node attributes to determine overall node health.",
                "unique":"0",
                "longdesc":"Requires external entities to create node attributes (named with the prefix '#health') with values: 'red', 'yellow' or 'green'.  Allowed values: none, migrate-on-red, only-green, progressive, custom"
            },
            "remove-after-stop":{
                "name":"remove-after-stop",
                "enabled":0,
                "value":"false",
                "content":{
                    "default":"false",
                    "type":"boolean"
                },
                "shortdesc":"Remove resources from the LRM after they are stopped",
                "unique":"0",
                "longdesc":"Always set this to false.  Other values are, at best, poorly tested and potentially dangerous."
            },
            "cluster-delay":{
                "name":"cluster-delay",
                "enabled":0,
                "value":"60s",
                "content":{
                    "default":"60s",
                    "type":"time"
                },
                "shortdesc":"Round trip delay over the network (excluding action execution)",
                "unique":"0",
                "longdesc":"The "correct" value will depend on the speed and load of your network and cluster nodes."
            }
        },
        "shortdesc":"Policy Engine Options",
        "version":"1.0",
        "nodecount":2,
        "isconfig":true,
        "longdesc":"This is a fake resource that details the options that can be configured for the Policy Engine."
    }
}
```



#### 2.1.2 修改集群属性

URI：/api/v1/haclusters/1

Method：PUT

请求

| 参数名称                | 是否必填          | 传入方式            | 参数类型           | 参数说明                 |
| :---------------------- | ---------------- | :----------------- | :----------------- | ----------------------- |
| no-quorum-policy        | 是               | json               | string             | 集群属性名及新的配置      |

请求示例：

```
{
    "no-quorum-policy": "stop"
}
```

响应格式：json
响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                         | string                             | 执行结果提示信息                       |

响应示例

```
{
    "action":true,
    "info":"Save crm metadata success"
}
```

#### 2.1.3 集群操作

URI：/api/v1/haclusters/1/:action

Method：DELETE

请求

| 参数名称        | 是否必填 | 传入方式 | 参数类型 | 参数说明   |
| :-------------- | -------- | :------- | :------- | ---------- |
| action              | 是       | path     | string   | 软件频道id |
| organization_id | 是       | path     | number   | 组织id     |

响应

| 参数名称                 | 参数类型 | 参数说明     |
| ------------------------ | -------- | ------------ |
| id                       | string   |              |
| label                    | string   |              |
| pending                  | Boolean  | 是否正在进行 |
| action                   | string   | 操作         |
| username                 | string   | 用户名       |
| started_at               | string   | 开始时间     |
| ended_at                 | string   | 结束时间     |
| state                    | string   |              |
| result                   | string   | 结果         |
| progress                 | number   |              |
| input                    | object   | 输入         |

响应示例

```
{
    "action":true,
    "info":"Save crm metadata success"
}
```
