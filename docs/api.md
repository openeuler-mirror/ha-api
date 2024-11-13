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

说明：获取集群属性

URI：/api/v1/haclusters/1

Method：GET

请求参数：无


响应格式：json

响应参数:

| 参数名称                          | 参数类型                               | 参数说明                             |
| :-------------------------------- | :------------------------------------ | ----------------------------------- |
| action                            | bool                                  | 返回结果状态                         |
| data                              | object                                | 返回的正常结果                       |
| error                             | string                                | 返回的错误信息                       |
| parameters                        | object                                | 各个组件的参数详情                   |
| shortdesc                         | string                                | 简短描述                             |
| version                           | string                                | 版本信息                             |
| nodecount                         | int                                   | 挂载的节点数量                       |
| isconfig                          | bool                                  | 是否是配置                           |
| longdesc                          | string                                | 完整描述                             |


响应示例:

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

说明：修改集群属性

URI：/api/v1/haclusters/1

Method：PUT

请求参数：

| 参数名称                   | 是否必填(其中一项必改） | 传入方式 | 参数类型    | 参数说明       |
|:-----------------------|--------------|:-----|:--------|------------|
| no-quorum-policy       | 是            | json | string  | 集群属性名及新的配置 |
| symmetric-cluster      | 是            | json | boolean | none       |
| maintenance-mode       | 是            | json | boolean | none       |
| start-failure-is-fatal | 是            | json | boolean | none       |
| stonith-enabled        | 是            | json | boolean | none       |
| node-health-strategy   | 是            | json | string  | none       |
| node-health-green      | 是            | json | string  | none       |
| node-health-yellow     | 是            | json | string  | none       |
| node-health-red        | 是            | json | string  | none       |

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


响应示例:

```
{
    "action":true,
    "info":"Save crm metadata success"
}
```

#### 2.1.3 集群操作

说明：集群操作

URI：/api/v1/haclusters/1/:action

Method：PUT

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                            |
| :---------------- | -------------- | :-------------- | :------------ | --------------------------------- |
| action            | 是             | path            | string        | 集群操作，start，stop或restart     |
| nodeid            | 是             | json            | string        | 节点id                             |
| nodeip            | 是             | json            | string        | 节点ip                             |
| password          | 是             | json            | string        | 节点password                       |

请求示例：

```
{
    "nodeauth":[
        {
            "nodeid": "ns187",
            "nodeip":"10.1.110.187",
            "password":"qwer1234"
        },
        {
            "nodeid": "ns188",
            "nodeip":"10.1.110.188",
            "password":"qwer1234"
        }
    ]
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 集群执行结果信息                      |

响应示例:

```
{
    "action":true,
    "info":"Save crm metadata success"`
}
```

##### 

### 2.2 节点

#### 2.2.1 获取节点列表

说明：获取节点列表

URI：/api/haclusters/1/nodes

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称   | 参数类型     | 参数说明                        |
|:-------|:---------|-----------------------------|
| action | bool     | 获取节点信息是否成功                  |
| data   | [object] | 节点信息，包含t以下的id，status以及is_dc |
| id     | string   | 节点名称                        |
| status | string   | 节点状态                        |
| is_dc  | string   | 是否为dc节点                     |

响应示例：

```
{
    "action":true,
    "data":[
        {
            "id":"ns187",
            "is_dc":true
            "status":"Master",
        },
        {
            "id":"ns188",
            "is_dc":false
            "status":"Not Running/Standby",
        }
    ]
}
```

#### 2.2.2 获取单个节点信息

说明：获取单个节点信息

URI：/api/haclusters/1/nodes/:nodeid

Method：GET

请求参数：

| 参数名称   | 是否必填 | 传入方式 | 参数类型   | 参数说明 |
|:-------|------|:-----|:-------|------|
| nodeid | 是    | path | string | 节点id |

响应格式：json

响应参数：

| 参数名称   | 参数类型   | 参数说明       |
|:-------|:-------|------------|
| action | bool   | 节点信息返回结果状态 |
| data   | string | 节点信息       |
| ips    | array  | 节点心跳ip列表   |

响应示例：

```
{   "action":true,
    "data": {
        'ips': ['10.1.110.188', '192.168.100.188']
    }
}
```


#### 2.2.3 节点操作

说明：节点操作

URI：/api/v1/haclusters/1/nodes/:node_id/:action

Method：PUT

请求参数：

| 参数名称    | 是否必填 | 传入方式 | 参数类型   | 参数说明                                             |
|:--------|------|:-----|:-------|--------------------------------------------------|
| node_id | 是    | path | string | 节点id                                             |
| action  | 是    | path | string | 支持修改的节点操作，包括unstandby、standby、stop、start和restart |

请求示例：

```
// start、stop和restart需要用户输入密码
{
    "password": "12345678",
}
```

响应格式：json

响应参数：

| 参数名称   | 参数类型   | 参数说明       |
|:-------|:-------|------------|
| action | bool   | 节点修改返回结果状态 |
| info   | string | 节点操作信息     |

响应示例：

```
{   "action":true,
    "info":"Change node status success"
}
```

### 2.3 资源

#### 2.3.1 获取资源列表

说明：获取资源列表

URI：/api/v1/haclusters/1/resources

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 资源信息列表                          |

响应示例：

```
{
    "action":true,
    "data":[
        {
            "status":"Running",
            "colocation":{
                "same_node":[
                    {
                        "with_rsc":"ms-drbd",
                        "rsc":"group-fs-ps"
                    }
                ],
                "same_node_num":1,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns187"
            ],
            "type":"group",
            "order":{
                "before_rscs":[
                    {
                        "id":"ms-drbd"
                    }
                ],
                "before_rscs_num":1,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "subrscs":[
                {
                    "status":"Running",
                    "running_node":[
                        "ns187"
                    ],
                    "type":"primitive",
                    "svc":"Filesystem",
                    "status_message":"",
                    "id":"fs-ps"
                }
            ],
            "status_message":"",
            "id":"group-fs-ps"
        },
        {
            "status":"Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns188"
            ],
            "type":"primitive",
            "svc":"Dummy",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "status_message":"",
            "id":"mydummy2"
        },
        {
            "status":"Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns188"
            ],
            "type":"group",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "subrscs":[
                {
                    "status":"Running",
                    "running_node":[
                        "ns188"
                    ],
                    "type":"primitive",
                    "svc":"mysql",
                    "status_message":"",
                    "id":"mysql"
                }
            ],
            "status_message":"",
            "id":"mysql-group"
        },
        {
            "status":"Not Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[

            ],
            "type":"primitive",
            "svc":"CTDB",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "status_message":"* test2_start_0 on ns187 'not installed' (5): call=70, status=complete, exitreason='Setup problem: couldn't find command: /usr/bin/tdbdump',* test2_start_0 on ns188 'not installed' (5): call=953, status=complete, exitreason='Setup problem: couldn't find command: /usr/bin/tdbdump',",
            "id":"test2"
        },
        {
            "status":"Not Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[

            ],
            "type":"group",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "subrscs":[
                {
                    "status":"Not Running",
                    "running_node":[

                    ],
                    "type":"primitive",
                    "svc":"CTDB",
                    "status_message":"",
                    "id":"test1"
                },
                {
                    "status":"Not Running",
                    "running_node":[

                    ],
                    "type":"primitive",
                    "svc":"Filesystem",
                    "status_message":"",
                    "id":"iscisi"
                }
            ],
            "status_message":"",
            "id":"group1"
        },
        {
            "status":"Running",
            "colocation":{
                "same_node":[
                    {
                        "with_rsc":"ms-drbd",
                        "rsc":"group-fs-ps"
                    }
                ],
                "same_node_num":1,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns187",
                "ns188"
            ],
            "type":"master",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[
                    {
                        "id":"group-fs-ps"
                    }
                ],
                "after_rscs_num":1
            },
            "location":[
                {
                    "node":"ns187",
                    "level":"Master Node"
                },
                {
                    "node":"ns188",
                    "level":"Slave 1"
                }
            ],
            "allow_unmigrate":false,
            "subrscs":[
                {
                    "status":"Running",
                    "running_node":[
                        "ns187"
                    ],
                    "type":"primitive",
                    "svc":"drbd",
                    "status_message":"",
                    "id":"drbd-ps:0"
                },
                {
                    "status":"Running",
                    "running_node":[
                        "ns188"
                    ],
                    "type":"primitive",
                    "svc":"drbd",
                    "status_message":"",
                    "id":"drbd-ps:1"
                }
            ],
            "status_message":"",
            "id":"ms-drbd"
        },
        {
            "status":"Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns188"
            ],
            "type":"primitive",
            "svc":"Dummy",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":true,
            "status_message":"",
            "id":"mydummy1"
        },
        {
            "status":"Not Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[

            ],
            "type":"primitive",
            "svc":"CTDB",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "status_message":"",
            "id":"y1"
        },
        {
            "status":"Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns187",
                "ns188"
            ],
            "type":"clone",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "subrscs":[
                {
                    "status":"Running",
                    "running_node":[
                        "ns187"
                    ],
                    "type":"primitive",
                    "svc":"IPaddr_6",
                    "status_message":"",
                    "id":"ip1:0"
                },
                {
                    "status":"Running",
                    "running_node":[
                        "ns188"
                    ],
                    "type":"primitive",
                    "svc":"IPaddr_6",
                    "status_message":"",
                    "id":"ip1:1"
                }
            ],
            "status_message":"",
            "id":"clone1"
        },
        {
            "status":"Not Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[

            ],
            "type":"primitive",
            "svc":"CTDB",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":false,
            "status_message":"",
            "id":"y2"
        },
        {
            "status":"Running",
            "colocation":{
                "same_node":[

                ],
                "same_node_num":0,
                "diff_node":[

                ],
                "diff_node_num":0
            },
            "running_node":[
                "ns188"
            ],
            "type":"group",
            "order":{
                "before_rscs":[

                ],
                "before_rscs_num":0,
                "after_rscs":[

                ],
                "after_rscs_num":0
            },
            "location":[

            ],
            "allow_unmigrate":true,
            "subrscs":[
                {
                    "status":"Running",
                    "running_node":[
                        "ns188"
                    ],
                    "type":"primitive",
                    "svc":"Dummy",
                    "status_message":"",
                    "id":"dummy4"
                }
            ],
            "status_message":"",
            "id":"dummy_group1"
        }
    ]
}
```


#### 2.3.2 添加资源

说明：添加资源

URI：/api/v1/haclusters/1/resources

Method：POST

请求参数：

| 参数名称                | 是否必填 | 传入方式 | 参数类型   | 参数说明   |
|:--------------------|------|:-----|:-------|--------|
| category            | 是    | json | string | 资源类型   |
| id                  | 是    | json | string | 资源名    |
| instance_attributes | 是    | json | object | none   |
| class               | 是    | json | string | 资源种类   |
| type                | 是    | json | string | 资源脚本名称 |
| provider            | 是    | json | string | 资源提供者  |

请求示例：

```
// primitive资源请求数据
{
    "category": "primitive",
    "meta_attributes":{
        "target-role":"Stopped"
    },
    "type":"CTDB", 
    "class":"ocf",
    "provider":"heartbeat",
    "instance_attributes":{
        "ctdb_recovery_lock":"lock"
    },
    "id":"test1"
}
// clone资源请求数据
{
  "category": "clone",    
    "id":"test5",
    "rsc_id":"test4",
    "meta_attributes":{
        "target-role":"Stopped"
    }
}
// group资源请求数据
{
     "category": "group",    
    "id":"tomcat_group",
    "rscs":[
              "tomcat6",
              "tomcat7"
    ],
    "meta_attributes":{
        "target-role":"Stopped"
    }
}
```

响应格式：json

响应参数：

| 参数名称                  | 参数类型                       | 参数说明                        |
| :------------------------ | :---------------------------- | ------------------------------- |
| action                    | bool                          | 返回结果状态                     |
| info                      | string                        | 返回结果信息                     |

响应示例:

```
{
    "action":True,
    'info':"Add primitive/clone/group resource success"
}
```

#### 2.3.3 单个资源操作

说明：单个资源操作

URI：/api/v1/haclusters/1/resources/:rsc_id/:action

Method：PUT

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                            |
| :---------------- | -------------- | :-------------- | :------------ | --------------------------------- |
| rsc_id            | 是             | path            | string        | 资源id                             |
| action            | 是             | path            | string        | 资源操作，包括start、stop、delete、cleanup、migrate、unmigrate、location、order、colocation等九种     |

请求示例：

```
// start、stop、delete和cleanup接口参数
{}
// migrate 接口参数
{
    "is_force": True,
    "to_node": "ns187",
    "period": "PYMDTHM3S"
}
// unmigrate 接口参数
{
    "rsc_id": "kk1",
    "is_all_rscs":False
}
// location 接口参数
{
    "node_level": [
        {
            "node": "ns187",
            "level": "Master Node"
        },
        {
            "node": "ns188",
            "level": "Slave 1"
        }
    ]
}
// colocation 接口参数
{
    "same_node": ["test1234"],
    "diff_node": ["group_tomcat"]
}
// order 口参数
{
    "before_rscs": ["test1234"],
    "after_rscs": ["group-fs-ps"]
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 资源操作信息                      |
| error                         | string                             | 资源操作失败信息                      |

响应示例：

```
{   "action":true,
    "info":"Action on resource success"
}
```


#### 2.3.4 获取所有资源创建数据

说明：获取所有资源创建数据

URI：/api/v1/haclusters/1/metas/:rsc_class/:rsc_type/:rsc_provider

Method：GET

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                            |
| :---------------- | -------------- | :-------------- | :------------ | --------------------------------- |
| rsc_class         | 是             | path            | string        | 资源类                             |
| rsc_type          | 是             | path            | string        | 资源类型                           |
| rsc_provider      | 是             | path            | string        | 资源provider                       |

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 资源创建数据信息                      |

响应示例：

```
{
    "action":true,
    "data":{
        "name":"CTDB",
        "parameters":[
            {
                "name":"ctdb_recovery_lock",
                "required":"1",
                "value":"",
                "content":{
                    "default":"",
                    "type":"string"
                },
                "shortdesc":"CTDB shared lock file",
                "unique":"1",
                "longdesc":"The location of a shared lock file, common across all nodes. This must be on shared storage, e.g.: /shared-fs/samba/ctdb.lock"
            },
            {
                "name":"ctdb_manages_samba",
                "required":"0",
                "value":"no",
                "content":{
                    "default":"no",
                    "type":"boolean"
                },
                "shortdesc":"Should CTDB manage Samba?",
                "unique":"0",
                "longdesc":"Should CTDB manage starting/stopping the Samba service for you? This will be deprecated in future, in favor of configuring a separate Samba resource."
            },
            {
                "name":"ctdb_manages_winbind",
                "required":"0",
                "value":"no",
                "content":{
                    "default":"no",
                    "type":"boolean"
                },
                "shortdesc":"Should CTDB manage Winbind?",
                "unique":"0",
                "longdesc":"Should CTDB manage starting/stopping the Winbind service for you? This will be deprecated in future, in favor of configuring a separate Winbind resource."
            },
            {
                "name":"ctdb_service_smb",
                "required":"0",
                "value":"",
                "content":{
                    "default":"",
                    "type":"string"
                },
                "shortdesc":"Name of smb init script",
                "unique":"0",
                "longdesc":"Name of smb init script.  Only necessary if CTDB is managing Samba directly.  Will usually be auto-detected."
            },
            {
                "name":"ctdb_service_nmb",
                "required":"0",
                "value":"",
                "content":{
                    "default":"",
                    "type":"string"
                },
                "shortdesc":"Name of nmb init script",
                "unique":"0",
                "longdesc":"Name of nmb init script.  Only necessary if CTDB is managing Samba directly.  Will usually be auto-detected."
            },
            {
                "name":"ctdb_service_winbind",
                "required":"0",
                "value":"",
                "content":{
                    "default":"",
                    "type":"string"
                },
                "shortdesc":"Name of winbind init script",
                "unique":"0",
                "longdesc":"Name of winbind init script.  Only necessary if CTDB is managing Winbind directly.  Will usually be auto-detected."
            },
            {
                "name":"ctdb_samba_skip_share_check",
                "required":"0",
                "value":"yes",
                "content":{
                    "default":"yes",
                    "type":"boolean"
                },
                "shortdesc":"Skip share check during monitor?",
                "unique":"0",
                "longdesc":"If there are very many shares it may not be feasible to check that all of them are available during each monitoring interval.  In that case this check can be disabled."
            },
            {
                "name":"ctdb_monitor_free_memory",
                "required":"0",
                "value":"100",
                "content":{
                    "default":"100",
                    "type":"integer"
                },
                "shortdesc":"Minimum amount of free memory (MB)",
                "unique":"0",
                "longdesc":"If the amount of free memory drops below this value the node will become unhealthy and ctdb and all managed services will be shutdown. Once this occurs, the administrator needs to find the reason for the OOM situation, rectify it and restart ctdb with "service ctdb start"."
            },
            {
                "name":"ctdb_start_as_disabled",
                "required":"0",
                "value":"no",
                "content":{
                    "default":"no",
                    "type":"boolean"
                },
                "shortdesc":"Start CTDB disabled?",
                "unique":"0",
                "longdesc":"When set to yes, the CTDB node will start in DISABLED mode and not host any public ip addresses."
            },
            {
                "name":"ctdb_config_dir",
                "required":"0",
                "value":"/etc/ctdb",
                "content":{
                    "default":"/etc/ctdb",
                    "type":"string"
                },
                "shortdesc":"CTDB config file directory",
                "unique":"0",
                "longdesc":"The directory containing various CTDB configuration files. The "nodes" and "notify.sh" scripts are expected to be in this directory, as is the "events.d" subdirectory."
            },
            {
                "name":"ctdb_binary",
                "required":"0",
                "value":"/usr/bin/ctdb",
                "content":{
                    "default":"/usr/bin/ctdb",
                    "type":"string"
                },
                "shortdesc":"CTDB binary path",
                "unique":"0",
                "longdesc":"Full path to the CTDB binary."
            },
            {
                "name":"ctdbd_binary",
                "required":"0",
                "value":"/usr/sbin/ctdbd",
                "content":{
                    "default":"/usr/sbin/ctdbd",
                    "type":"string"
                },
                "shortdesc":"CTDB Daemon binary path",
                "unique":"0",
                "longdesc":"Full path to the CTDB cluster daemon binary."
            },
            {
                "name":"ctdb_socket",
                "required":"0",
                "value":"/run/ctdb/ctdbd.socket",
                "content":{
                    "default":"/run/ctdb/ctdbd.socket",
                    "type":"string"
                },
                "shortdesc":"CTDB socket location",
                "unique":"1",
                "longdesc":"Full path to the domain socket that ctdbd will create, used for local clients to attach and communicate with the ctdb daemon."
            },
            {
                "name":"ctdb_dbdir",
                "required":"0",
                "value":"/var/run",
                "content":{
                    "default":"/var/run",
                    "type":"string"
                },
                "shortdesc":"CTDB database directory",
                "unique":"1",
                "longdesc":"The directory to put the local CTDB database files in. Persistent database files will be put in ctdb_dbdir/persistent."
            },
            {
                "name":"ctdb_logfile",
                "required":"0",
                "value":"/var/log/ctdb/log.ctdb",
                "content":{
                    "default":"/var/log/ctdb/log.ctdb",
                    "type":"string"
                },
                "shortdesc":"CTDB log file location",
                "unique":"0",
                "longdesc":"Full path to log file. To log to syslog instead, use the value "syslog"."
            },
            {
                "name":"ctdb_rundir",
                "required":"0",
                "value":"/run/ctdb",
                "content":{
                    "default":"/run/ctdb",
                    "type":"string"
                },
                "shortdesc":"CTDB runtime directory location",
                "unique":"0",
                "longdesc":"Full path to ctdb runtime directory, used for storage of socket lock state."
            },
            {
                "name":"ctdb_debuglevel",
                "required":"0",
                "value":"2",
                "content":{
                    "default":"2",
                    "type":"integer"
                },
                "shortdesc":"CTDB debug level",
                "unique":"0",
                "longdesc":"What debug level to run at (0-10). Higher means more verbose."
            },
            {
                "name":"smb_conf",
                "required":"0",
                "value":"/etc/samba/smb.conf",
                "content":{
                    "default":"/etc/samba/smb.conf",
                    "type":"string"
                },
                "shortdesc":"Path to smb.conf",
                "unique":"0",
                "longdesc":"Path to default samba config file.  Only necessary if CTDB is managing Samba."
            },
            {
                "name":"smb_private_dir",
                "required":"0",
                "value":"",
                "content":{
                    "default":"",
                    "type":"string"
                },
                "shortdesc":"Samba private dir (deprecated)",
                "unique":"1",
                "longdesc":"The directory for smbd to use for storing such files as smbpasswd and secrets.tdb.  Old versions of CTBD (prior to 1.0.50) required this to be on shared storage.  This parameter should not be set for current versions of CTDB, and only remains in the RA for backwards compatibility."
            },
            {
                "name":"smb_passdb_backend",
                "required":"0",
                "value":"tdbsam",
                "content":{
                    "default":"tdbsam",
                    "type":"string"
                },
                "shortdesc":"Samba passdb backend",
                "unique":"0",
                "longdesc":"Which backend to use for storing user and possibly group information.  Only necessary if CTDB is managing Samba."
            },
            {
                "name":"smb_idmap_backend",
                "required":"0",
                "value":"tdb2",
                "content":{
                    "default":"tdb2",
                    "type":"string"
                },
                "shortdesc":"Samba idmap backend",
                "unique":"0",
                "longdesc":"Which backend to use for SID/uid/gid mapping.  Only necessary if CTDB is managing Samba."
            },
            {
                "name":"smb_fileid_algorithm",
                "required":"0",
                "value":"",
                "content":{
                    "default":"",
                    "type":"string"
                },
                "shortdesc":"Samba VFS fileid algorithm",
                "unique":"0",
                "longdesc":"Which fileid:algorithm to use with vfs_fileid.  The correct value depends on which clustered filesystem is in use, e.g.:for OCFS2, this should be set to "fsid".  Only necessary if CTDB is managing Samba."
            }
        ],
        "actions":[
            {
                "interval":"0",
                "name":"start",
                "timeout":"90"
            },
            {
                "interval":"0",
                "name":"stop",
                "timeout":"100"
            },
            {
                "depth":"0",
                "interval":"10",
                "name":"monitor",
                "timeout":"20"
            },
            {
                "interval":"0",
                "name":"meta-data",
                "timeout":"5"
            },
            {
                "interval":"0",
                "name":"validate-all",
                "timeout":"30"
            }
        ],
        "version":"1.0",
        "shortdesc":"CTDB Resource Agent",
        "longdesc":"This resource agent manages CTDB, allowing one to use Clustered Samba in a Linux-HA/Pacemaker cluster.  You need a shared filesystem (e.g. OCFS2 or GFS2) on which the CTDB lock will be stored.  Create /etc/ctdb/nodes containing a list of private IP addresses of each node in the cluster, then configure this RA as a clone.  This agent expects the samba and windbind resources to be managed outside of CTDB's control as a separate set of resources controlled by the cluster manager.  The optional support for enabling CTDB management of these daemons will be depreciated. For more information see http://linux-ha.org/wiki/CTDB_(resource_agent)"
    }
}
```

#### 2.3.5 获取所有资源类型

说明：获取所有资源类型

URI：/api/v1/haclusters/1/metas

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 各种资源类型清单                      |

响应示例：

```
{
    "action":true,
    "data":{
        "stonith":[
            "apcmaster",
            "apcmastersnmp",
            "apcsmart",
            "baytech",
            "bladehpi",
            "cyclades",
            "fence_apc",
            "fence_apc_snmp",
            "fence_bladecenter",
            "fence_brocade",
            "fence_cisco_mds",
            "fence_cisco_ucs",
            "fence_compute",
            "fence_drac5",
            "fence_eaton_snmp",
            "fence_emerson",
            "fence_eps",
            "fence_hpblade",
            "fence_ibmblade",
            "fence_idrac",
            "fence_ifmib",
            "fence_ilo",
            "fence_ilo2",
            "fence_ilo3",
            "fence_ilo3_ssh",
            "fence_ilo4",
            "fence_ilo4_ssh",
            "fence_ilo_moonshot",
            "fence_ilo_mp",
            "fence_ilo_ssh",
            "fence_imm",
            "fence_intelmodular",
            "fence_ipdu",
            "fence_ipmilan",
            "fence_kdump",
            "fence_mpath",
            "fence_rhevm",
            "fence_rsa",
            "fence_rsb",
            "fence_sbd",
            "fence_scsi",
            "fence_virsh",
            "fence_virt",
            "fence_vmware_soap",
            "fence_wti",
            "fence_xvm",
            "ibmhmc",
            "ipmilan",
            "meatware",
            "nw_rpc100s",
            "rcd_serial",
            "rps10",
            "suicide",
            "wti_mpc",
            "wti_nps"
        ],
        "systemd":[
            "ModemManager",
            "NetworkManager",
            "NetworkManager-wait-online",
            "abrt-ccpp",
            "abrt-oops",
            "abrt-vmcore",
            "abrt-xorg",
            "abrtd",
            "accounts-daemon",
            "alsa-restore",
            "alsa-state",
            "alsa-store",
            "apparmor",
            "atd",
            "auditd",
            "auth-rpcgss-module",
            "avahi-daemon",
            "blk-availability",
            "brandbot",
            "chronyd",
            "cobra",
            "colord",
            "corosync",
            "cpupower",
            "crond",
            "cups",
            "dbus",
            "dm-event",
            "dmraid-activation",
            "dracut-shutdown",
            "emergency",
            "exim",
            "fcoe",
            "gdm",
            "getty@tty1",
            "getty@tty2",
            "gssproxy",
            "ha-api",
            "hypervfcopyd",
            "hypervkvpd",
            "hypervvssd",
            "icinga",
            "ido2db",
            "ip6tables",
            "iptables",
            "irqbalance",
            "iscsi",
            "iscsi-shutdown",
            "iscsid",
            "iscsiuio",
            "kdump",
            "kmod-static-nodes",
            "ksm",
            "ksmtuned",
            "ldconfig",
            "libstoragemgmt",
            "libvirt-guests",
            "libvirtd",
            "livesys",
            "lldpad",
            "lvm2-activation",
            "lvm2-activation-early",
            "lvm2-lvmetad",
            "lvm2-lvmpolld",
            "lvm2-monitor",
            "lvm2-pvscan@252:2",
            "mdmonitor",
            "microcode",
            "multipathd",
            "mysqld",
            "neokylinhautils",
            "network",
            "nfs-config",
            "nfs-idmapd",
            "nfs-mountd",
            "nfs-server",
            "nfs-utils",
            "nkucsd",
            "ntpd",
            "ntpdate",
            "pacemaker",
            "pacemaker-mgmt",
            "packagekit",
            "plymouth-quit",
            "plymouth-quit-wait",
            "plymouth-read-write",
            "plymouth-start",
            "polkit",
            "postfix",
            "postgresql",
            "qemu-guest-agent",
            "rc-local",
            "rescue",
            "rhel-autorelabel",
            "rhel-autorelabel-mark",
            "rhel-configure",
            "rhel-dmesg",
            "rhel-import-state",
            "rhel-loadmodules",
            "rhel-readonly",
            "rngd",
            "rpc-gssd",
            "rpc-statd",
            "rpc-statd-notify",
            "rpc-svcgssd",
            "rpcbind",
            "rsyslog",
            "rtkit-daemon",
            "sendmail",
            "sm-client",
            "smartd",
            "sntp",
            "spice-vdagentd",
            "sshd",
            "sshd-keygen",
            "syslog",
            "sysnotify",
            "sysstat",
            "systemcenter",
            "systemd-ask-password-console",
            "systemd-ask-password-plymouth",
            "systemd-ask-password-wall",
            "systemd-binfmt",
            "systemd-firstboot",
            "systemd-fsck-root",
            "systemd-hwdb-update",
            "systemd-initctl",
            "systemd-journal-catalog-update",
            "systemd-journal-flush",
            "systemd-journald",
            "systemd-logind",
            "systemd-machine-id-commit",
            "systemd-modules-load",
            "systemd-random-seed",
            "systemd-random-seed-load",
            "systemd-readahead-collect",
            "systemd-readahead-done",
            "systemd-readahead-replay",
            "systemd-reboot",
            "systemd-remount-fs",
            "systemd-shutdownd",
            "systemd-sysctl",
            "systemd-sysusers",
            "systemd-tmpfiles-clean",
            "systemd-tmpfiles-setup",
            "systemd-tmpfiles-setup-dev",
            "systemd-udev-settle",
            "systemd-udev-trigger",
            "systemd-udevd",
            "systemd-update-done",
            "systemd-update-utmp",
            "systemd-update-utmp-runlevel",
            "systemd-user-sessions",
            "systemd-vconsole-setup",
            "tuned",
            "udisks2",
            "unbound-anchor",
            "upower",
            "vmtoolsd",
            "wpa_supplicant",
            "ypbind"
        ],
        "lsb":[
            "CobraApi",
            "cobra",
            "cobrastatus",
            "ha-api",
            "icinga",
            "ido2db",
            "netconsole",
            "network",
            "nkucsd",
            "npcd",
            "nrpe",
            "systemcenter"
        ],
        "service":[
            "CobraApi",
            "ModemManager",
            "NetworkManager",
            "NetworkManager-wait-online",
            "abrt-ccpp",
            "abrt-oops",
            "abrt-vmcore",
            "abrt-xorg",
            "abrtd",
            "accounts-daemon",
            "alsa-restore",
            "alsa-state",
            "alsa-store",
            "apparmor",
            "atd",
            "auditd",
            "auth-rpcgss-module",
            "avahi-daemon",
            "blk-availability",
            "brandbot",
            "chronyd",
            "cobra",
            "cobra",
            "cobrastatus",
            "colord",
            "corosync",
            "cpupower",
            "crond",
            "cups",
            "dbus",
            "dm-event",
            "dmraid-activation",
            "dracut-shutdown",
            "emergency",
            "exim",
            "fcoe",
            "gdm",
            "getty@tty1",
            "getty@tty2",
            "gssproxy",
            "ha-api",
            "ha-api",
            "hypervfcopyd",
            "hypervkvpd",
            "hypervvssd",
            "icinga",
            "icinga",
            "ido2db",
            "ido2db",
            "ip6tables",
            "iptables",
            "irqbalance",
            "iscsi",
            "iscsi-shutdown",
            "iscsid",
            "iscsiuio",
            "kdump",
            "kmod-static-nodes",
            "ksm",
            "ksmtuned",
            "ldconfig",
            "libstoragemgmt",
            "libvirt-guests",
            "libvirtd",
            "livesys",
            "lldpad",
            "lvm2-activation",
            "lvm2-activation-early",
            "lvm2-lvmetad",
            "lvm2-lvmpolld",
            "lvm2-monitor",
            "lvm2-pvscan@252:2",
            "mdmonitor",
            "microcode",
            "multipathd",
            "mysqld",
            "neokylinhautils",
            "netconsole",
            "network",
            "network",
            "nfs-config",
            "nfs-idmapd",
            "nfs-mountd",
            "nfs-server",
            "nfs-utils",
            "nkucsd",
            "nkucsd",
            "npcd",
            "nrpe",
            "ntpd",
            "ntpdate",
            "pacemaker",
            "pacemaker-mgmt",
            "packagekit",
            "plymouth-quit",
            "plymouth-quit-wait",
            "plymouth-read-write",
            "plymouth-start",
            "polkit",
            "postfix",
            "postgresql",
            "qemu-guest-agent",
            "rc-local",
            "rescue",
            "rhel-autorelabel",
            "rhel-autorelabel-mark",
            "rhel-configure",
            "rhel-dmesg",
            "rhel-import-state",
            "rhel-loadmodules",
            "rhel-readonly",
            "rngd",
            "rpc-gssd",
            "rpc-statd",
            "rpc-statd-notify",
            "rpc-svcgssd",
            "rpcbind",
            "rsyslog",
            "rtkit-daemon",
            "sendmail",
            "sm-client",
            "smartd",
            "sntp",
            "spice-vdagentd",
            "sshd",
            "sshd-keygen",
            "syslog",
            "sysnotify",
            "sysstat",
            "systemcenter",
            "systemcenter",
            "systemd-ask-password-console",
            "systemd-ask-password-plymouth",
            "systemd-ask-password-wall",
            "systemd-binfmt",
            "systemd-firstboot",
            "systemd-fsck-root",
            "systemd-hwdb-update",
            "systemd-initctl",
            "systemd-journal-catalog-update",
            "systemd-journal-flush",
            "systemd-journald",
            "systemd-logind",
            "systemd-machine-id-commit",
            "systemd-modules-load",
            "systemd-random-seed",
            "systemd-random-seed-load",
            "systemd-readahead-collect",
            "systemd-readahead-done",
            "systemd-readahead-replay",
            "systemd-reboot",
            "systemd-remount-fs",
            "systemd-shutdownd",
            "systemd-sysctl",
            "systemd-sysusers",
            "systemd-tmpfiles-clean",
            "systemd-tmpfiles-setup",
            "systemd-tmpfiles-setup-dev",
            "systemd-udev-settle",
            "systemd-udev-trigger",
            "systemd-udevd",
            "systemd-update-done",
            "systemd-update-utmp",
            "systemd-update-utmp-runlevel",
            "systemd-user-sessions",
            "systemd-vconsole-setup",
            "tuned",
            "udisks2",
            "unbound-anchor",
            "upower",
            "vmtoolsd",
            "wpa_supplicant",
            "ypbind"
        ],
        "ocf":{
            "heartbeat":[
                "CTDB",
                "Cpu_Health",
                "Delay",
                "Disk_Health",
                "Filesystem",
                "IPaddr",
                "IPaddr2",
                "IPaddr_6",
                "IPsrcaddr",
                "IPv6addr",
                "LVM",
                "MailTo",
                "Mem_Health",
                "Route",
                "SAPDatabase",
                "SAPHana",
                "SAPHanaTopology",
                "SAPInstance",
                "SendArp",
                "Squid",
                "VirtualDomain",
                "Xinetd",
                "apache",
                "clvm",
                "conntrackd",
                "db2",
                "dhcpd",
                "docker",
                "ethmonitor",
                "exportfs",
                "galera",
                "iSCSILogicalUnit",
                "iSCSITarget",
                "iface-vlan",
                "mysql",
                "named",
                "nfsnotify",
                "nfsserver",
                "nginx",
                "oracle",
                "oralsnr",
                "pgsql",
                "postfix",
                "rabbitmq-cluster",
                "redis",
                "rsyncd",
                "slapd",
                "symlink",
                "tomcat"
            ],
            "linbit":[
                "drbd"
            ],
            "pacemaker":[
                "ClusterMon",
                "Dummy",
                "HealthCPU",
                "HealthSMART",
                "Stateful",
                "SysInfo",
                "SystemHealth",
                "controld",
                "ping",
                "pingd",
                "remote"
            ],
            "openstack":[
                "NovaCompute",
                "NovaEvacuate"
            ]
        }
    }
}
```

#### 2.3.6 获取资源元属性

说明：获取资源元属性

URI：/api/v1/haclusters/1/resources/meta_attributes/:catagory

Method：GET

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                                      |
| :--------------- | -------------- | :-------------- | :------------ | -------------------------------------------- |
| catagory         | 是             | path            | string        | 资源类型，clone、primitive或者group           |

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 资源元信息                            |

响应示例:

```
{
    "action":true,
    "data":{
        "migration-threshold":{
            "content":{
                "type":"integer"
            },
            "name":"migration-threshold"
        },
        "allow-migrate":{
            "content":{
                "type":"boolean"
            },
            "name":"allow-migrate"
        },
        "failure-timeout":{
            "content":{
                "type":"integer"
            },
            "name":"failure-timeout"
        },
        "priority":{
            "content":{
                "type":"integer"
            },
            "name":"priority"
        },
        "resource-stickiness":{
            "content":{
                "type":"integer"
            },
            "name":"resource-stickiness"
        },
        "target-role":{
            "content":{
                "default":"Stopped",
                "type":"enum",
                "values":[
                    "Stopped",
                    "Started"
                ]
            },
            "name":"target-role"
        },
        "multiple-active":{
            "content":{
                "type":"enum",
                "values":[
                    "stop_start",
                    "stop_only",
                    "block"
                ]
            },
            "name":"multiple-active"
        },
        "is-managed":{
            "content":{
                "type":"boolean"
            },
            "name":"is-managed"
        }
    }
}
```

#### 2.3.7 获取relation信息

说明：获取relation信息

URI：/api/v1/haclusters/1/resources/:rsc_id/relations/:relation

Method：GET

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                                      |
| :--------------- | -------------- | :-------------- | :------------ | -------------------------------------------- |
| rsc_id           | 是             | path            | string        | 资源id                                        |
| relation         | 是             | path            | string        | relation类型，location、order或者colocation   |

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | relation配置信息                          |

响应示例:

```
// location返回结果
{
    "action":true,
    "data":{
        "node_level":[
            {
                "node":"ns187",
                "level":"Master Node"
            },
            {
                "node":"ns188",
                "level":"Slave 1"
            }
        ],
        "rsc_id":"iscisi"
    }
}

// colocation返回结果
{
    "action":true,
    "data":{
        "same_node":[

        ],
        "rsc_id":"iscisi",
        "diff_node":[
            "dummy_group1"
        ]
    }
}
// order返回结果
{
    "action":true,
    "data":{
        "before_rscs":[

        ],
        "rsc_id":"iscisi",
        "after_rscs":[
            "dummy_group1"
        ]
    }
}
```

#### 2.3.8 修改资源

说明：修改资源

URI：/api/v1/haclusters/1/resources/:rsc_id

Method：PUT

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                   |
| :--------------- | -------------- | :-------------- | :------------ | -------------------------- |
| rsc_id           | 是             | path            | string        | 资源id                     |

请求示例：

```
// 修改primitive资源
{
    "category": "primitive",
    "actions":[
        {
            "interval":"100",
            "name":"start"
        }
    ],
    "meta_attributes":{
        "resource-stickiness":"104",
        "is-managed":"true",
        "target-role":"Started"
    },
    "type":"Filesystem",
    "id":"iscisi",
    "provider":"heartbeat",
    "instance_attributes":{
        "device":"/dev/sda1",
        "directory":"/var/lib/mysql",
        "fstype":"ext4"
    },
    "class":"ocf"
}
// 修改group资源
{
    "id":"group1",
    "category":"group",
    "rscs":["iscisi", "test1" ],
    "meta_attributes":{
        "target-role":"Started"
    }
}
// 修改clone资源
{
    "id":"clone1",
    "category":"clone",
    "rsc_id":"ip1",
    "meta_attributes":{
        "target-role":"Stopped"
    }
}
```


响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | object                             | 资源配置信息                          |

响应示例:

```
// 修改primitive资源
{
    'action':True,
    'info':"Set rsc info success"
}
// 修改group资源
{
    'action':True,
    'info':"Set rsc info success"
}
// 修改clone资源
{
    'action':True,
    'info':"Set rsc info success"
}
```

#### 2.3.9 获取资源

说明：获取资源信息

URI：/api/v1/haclusters/1/resources/:rsc_id

Method：GET

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                   |
| :--------------- | -------------- | :-------------- | :------------ | -------------------------- |
| rsc_id           | 是             | path            | string        | 资源id                     |

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 资源配置信息                          |

响应示例:

```
// 获取primitive资源返回结果
{
    "action":true,
    "data":{
        "category":"primitive",
        "meta_attributes":{
            "resource-stickiness":"104"
        },
        "instance_attributes":{
            "device":"/dev/sda1",
            "directory":"/var/lib/mysql",
            "fstype":"ext4"
        },
        "actions":[
            {
                "interval":"103",
                "timeout":"61",
                "name":"start"
            },
            {
                "interval":"23",
                "timeout":"41",
                "name":"monitor"
            }
        ],
        "id":"iscisi",
        "provider":"heartbeat",
        "type":"Filesystem",
        "class":"ocf"
    }
}
// 获取group资源返回结果
{
    "action":true,
    "data":{
        "category":"group",
        "rscs":[
            "iscisi",
            "test1"
        ],
        "id":"group1",
        "meta_attributes":{
            "target-role":"Stopped"
        }
    }
}
// 获取clone资源返回结果
{
    "action":true,
    "data":{
        "category":"clone",
        "rsc_id":"ip1",
        "id":"clone1",
        "meta_attributes":{
            "target-role":"Started"
        }
    }
}
```

### 2.4 登录

#### 2.4.1 登录接口

说明：用户登录

URI：/api/v1/login

Method：POST

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型       | 参数说明                   |
| :--------------- | -------------- | :-------------- | :------------ | -------------------------- |
| username         | 是             | json            | string        | 用户名                     |
| password         | 是             | json            | string        | 密码                       |

请求示例：

```
{
    "username": "hacluster",
    "password": "12345678"
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |

响应示例:

```
{
    "action": true
}
```

### 2.5 告警

#### 2.5.1 获取告警信息

说明：获取告警信息

URI：/api/v1/haclusters/1/alarms

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | int                                | 返回数据                             |
| sender                        | string                             | 发送邮箱                             |
| smtp                          | string                             | 邮箱smtp                             |
| flag                          | bool                               | switch标志                           |
| receiver                      | array                              | 接收者邮箱列表                        |
| password                      | string                             | 邮箱密码                             |
| port                          | string                             | 邮箱端口                             |

响应示例:

```
{
    "action":true,
    "data":{
        "sender":"hatest@cs2c.com.cn",
        "smtp":"mail.cs2c.com.cn",
        "flag":true,
        "receiver":[
            "abc1@cs2c.com.cn",
            "abc@cs2c.com.cn"
        ],
        "password":"",
        "port":"25"
    }
}
```

#### 2.5.1 创建告警

说明：创建告警

URI：/api/v1/haclusters/1/alarms

Method：POST

请求参数：

| 参数名称          | 是否必填       | 传入方式        | 参数类型      | 参数说明                   |
| :--------------- | ------------- | :------------- | :------------ | -------------------------- |
| sender           | 是            | json           | string        | 发送邮箱                    |
| smtp             | 双            | json           | string        | 邮箱smtp                    |
| flag             | 是            | json           | bool          | switch标志                  |
| receiver         | 是            | json           | array         | 接收者邮箱列表               |
| password         | 是            | json           | string        | 邮箱密码                    |
| port             | 是            | json           | string        | 邮箱端口                    |

请求示例：

```
{
    "sender":"hatest@cs2c.com.cn",
    "smtp":"mail.cs2c.com.cn",
    "flag":true,
    "receiver":[
        "abc1@cs2c.com.cn",
        "abc@cs2c.com.cn"
    ],
    "password":"BFRBBQFaVFRcREFHUCkaHw==",
    "port":"25"
}
```

响应格式：json

响应参数：

| 参数名称                  | 参数类型                       | 参数说明                        |
| :------------------------ | :---------------------------- | ------------------------------- |
| action                    | bool                          | 返回结果状态                     |
| info                      | string                        | 返回结果信息                     |

响应示例:

```
{
    "action":True,
    'info':"Set alarm success"
}
```

#### 2.5.1 删除告警信息

TODO:

### 2.6 心跳

#### 2.6.1 获取网络心跳信息

说明：获取网络心跳信息

URI：/api/v1/haclusters/1/configs

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | int                                | 返回心跳状态，0：正常，非0：异常       |
| hbaddrs1                      | object                             | 单心跳节点配置                        |
| hbaddrs2                      | object                             | 双心跳节点配置                        |
| ip                            | string                             | 节点ip                               |
| nodeid                        | string                             | 节点id名                             |
| hbaddrs2_enabled              | int                                | 双心跳是否开启                        |

响应示例:

```
{
    "action":true,
    "data":{        
        "hbaddrs1":[
            {
                "ip":"192.168.100.187",
                "nodeid": "ns187"               
            },
            {
                "ip":"192.168.100.188",
                "nodeid": "ns188"
            }
        ],
        "hbaddrs2":[
            {
                "ip":"192.168.100.187",
                "nodeid": "ns187"
            },
            {
                "ip":"192.168.100.188",
                "nodeid": "ns188"
            }
        ],
        "hbaddrs2_enabled": 1
	}
}
```

#### 2.6.2 创建/编辑心跳信息

说明：创建/修改心跳信息

URI：/api/v1/haclusters/1/configs

Method：POST

请求参数：

| 参数名称             | 是否必填       | 传入方式        | 参数类型      | 参数说明                   |
| :------------------- | ------------- | :------------- | :------------ | -------------------------- |
| hbaddrs1             | 是            | json           | object        | 单心跳节点配置              |
| hbaddrs2             | 双            | json           | object        | 双心跳节点配置              |
| ip                   | 是            | json           | string        | 节点ip                     |
| nodeid               | 是            | json           | string        | 节点id名                   |
| hbaddrs2_enabled     | 是            | json           | int           | 双心跳是否开启              |

请求示例：

```
{        
    "hbaddrs1":[
        {
            "ip":"192.168.100.187",
            "nodeid": "ns187"
        },
        {
            "ip":"192.168.100.188",
            "nodeid": "ns188"
        }
    ],

    "hbaddrs2":[
        {
            "ip":"192.168.100.187",
            "nodeid": "ns187"
        },
        {
            "ip":"192.168.100.188",
            "nodeid": "ns188"
        }
    ],
    "hbaddrs2_enabled":  1
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 返回结果信息                          |

响应示例:

```
{
    "action":True,
    'info':"update success"
}
```

#### 2.6.3 获取可用磁盘心跳设备列表

TODO: 

#### 2.6.4 获取心跳状态信息

说明：获取心跳状态信息

URI：/api/v1/haclusters/1/hbstatus

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | int                                | 返回心跳状态，0：正常，非0：异常       |

响应示例:

```
{
    "action":True,
    "data": 0
}
```


### 2.7 DRBD

#### 2.7.1 获取DRBD信息

说明：获取DRBD信息

URI：/api/v1/haclusters/1/drbd

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 返回结果信息                          |
| dir                           | string                             | 路径                                 |
| main_node_ip                  | string                             | 主节点ip                             |
| main_node_name                | string                             | 主节点名称                           |
| main_node_device              | string                             | 主节点设备                           |
| slave_node_ip                 | string                             | 从节点ip                             |
| slave_node_name               | string                             | 从节点名称                           |
| slave_node_device             | string                             | 从节点设备                           |
| drbd_type                     | string                             | drbd类型，0：主从，1：主主            |

响应示例:

```
{
    "action": true,
    "data":{
        "main_node_name": "ns187",
        "dir": "/drbd",
        "main_node_device":"/dev/vdb",
        "slave_node_ip": "10.1.110.188",
        "slave_node_name":"ns188",
        "main_node_ip":"10.1.110.187",
        "slave_node_device":"/dev/vdb",
        "drbd_type":  "1"
    }
}
```

#### 2.7.2 创建及修改DRBD

说明：创建/修改DRBD

URI：/api/v1/haclusters/1/drbd

Method：POST

请求参数：

| 参数名称             | 是否必填       | 传入方式        | 参数类型      | 参数说明                   |
| :------------------- | ------------- | :------------- | :------------ | -------------------------- |
| dir                  | 是            | json           | string        | 路径                        |
| main_node_ip         | 是            | json           | string        | 主节点ip                    |
| main_node_name       | 是            | json           | string        | 主节点名称                  |
| main_node_device     | 是            | json           | string        | 主节点设备                  |
| slave_node_ip        | 是            | json           | string        | 从节点ip                    |
| slave_node_name      | 是            | json           | string        | 从节点名称                  |
| slave_node_device    | 是            | json           | string        | 从节点设备                  |
| drbd_type            | 是            | json           | string        | drbd类型，0：主从，1：主主   |

请求示例：

```
{
    "main_node_name":"ns187",
    "dir":"/drbd",
    "main_node_device":"/dev/vdb",
    "slave_node_ip":"10.1.110.188",
    "slave_node_name":"ns188",
    "main_node_ip":"10.1.110.187",
    "slave_node_device":"/dev/vdb",
    "drbd_type":  "1"
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 返回结果信息                          |

响应示例:

```
{
    "action":True,
    'info':"Save drbd configuration success"
}
```

#### 2.7.3 获取DRBD状态

TODO: 

#### 2.7.4 删除DRBD

说明：删除DRBD

URI：/api/v1/haclusters/1/drbd

Method：DELETE

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 返回结果信息                          |

响应示例:

```
{
    "action":True,
    'info':"Delete drbd configuration success"
}
```

### 2.8 GFS

#### 2.8.1 获取GFS信息

说明：获取GFS资源信息

URI：/api/v1/haclusters/1/gfs

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 返回结果信息                          |
| device                        | string                             | 设备路径                              |
| max_num                       | int                                | 最大数量                              |
| mount_dir                     | string                             | 设备挂载路径                          |

响应示例:

```
{
    "action": true,
    "data":{
        "device": "/dev/vdb",
        "max_num": 2,
        "mount_dir": "/drbd"
    }
}
```


#### 2.8.2 创建及修改GFS

说明：创建/修改GFS资源

URI：/api/v1/haclusters/1/gfs

Method：POST

请求参数：

| 参数名称           | 是否必填          | 传入方式            | 参数类型           | 参数说明                  |
| :----------------- | ---------------- | :----------------- | :----------------- | ------------------------ |
| device             | 是               | json               | string             | 设备路径                  |
| max_num            | 是               | json               | int                | 最大数量                  |
| mount_dir          | 是               | json               | string             | 设备挂载路径              |

请求示例：

```
{
    "device": "/dev/vdb",
    "max_num": 2,
    "mount_dir": "/drbd"
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 返回结果信息                          |

响应示例:

```
{
    "action":True,
    'info':"Save gfs configuration success"
}
```

#### 2.8.3 删除GFS

说明：删除挂载的GFS

URI：/api/v1/haclusters/1/gfs

Method：DELETE

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 返回结果信息                          |

响应示例:

```
{
    "action":True,
    'info':"Delete gfs success"
}
```

### 2.9 脚本

#### 2.9.1 生成脚本

说明：生成资源控制脚本

URI：/api/v1/haclusters/1/scripts

Method：POST

请求参数：

| 参数名称           | 是否必填          | 传入方式            | 参数类型           | 参数说明                  |
| :----------------- | ---------------- | :----------------- | :----------------- | ------------------------ |
| name               | 是               | json               | string             | 脚本名称                  |
| start              | 是               | json               | string             | start命令脚本             |
| stop               | 是               | json               | string             | stop命令脚本              |
| monitor            | 是               | json               | string             | monitor命令脚本           |

请求示例：

```
{
    "name": "tomcat3",
    "start":"/usr/local/bin/tomcat start",
    "stop": "/usr/local/bin/tomcat stop",
    "monitor": "/usr/local/bin/tomcat monitor"
}
```

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 返回结果信息                          |

响应示例:

```
{
    "action": true,
    "info": "Create script success!"
}
```


### 2.10 日志

#### 2.10.1 生成日志

说明：生成集群日志并返回日志路径

URI：/api/v1/haclusters/1/logs

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 返回结果数据                          |
| filepath                      | string                             | 生成的日志压缩文件的路径               |

响应示例:

```
{
    'action': True,
    'data': {
        "filepath": "static/neokylinha-log.tar"
    }
}
```

### 2.11 指令集

#### 2.11.1 获取指令执行结果

说明：执行预定义的指令，返回执行结果

URI：/api/v1/haclusters/1/commands/:cmd_type

Method：GET

请求参数：

| 参数名称           | 是否必填          | 传入方式            | 参数类型           | 参数说明                 |
| :----------------- | ---------------- | :----------------- | :----------------- | ------------------------ |
| cmd_type           | 是               | path               | int                | 预定义的指令类型          |

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | string                             | 预定义的指令执行结果                  |

响应示例:

```
{
    "action": true,
    "data": "Stack: corosync\nCurrent DC: ns175 (version 1.1.16-12.el7.1-94ff4df) - partition with quorum\nLast updated: Thu Jun  7 15:55:05 2018\nLast change: Wed Jun  6 10:56:10 2018 by root via mgmtd on ns175\n\n2 nodes configured\n19 resources configured (23 DISABLED)\n\nOnline: [ ns175 ]\nOFFLINE: [ ns176 ]\n\nActive resources:\n\n dummy\t(ocf::heartbeat:Dummy):\tStarted ns175\n test2\t(ocf::heartbeat:Dummy):\tStarted ns175\n\nOperations:\n* Node ns175:\n   dummy: migration-threshold=1000000\n    + (60) start: rc=0 (ok)\n   test2: migration-threshold=1000000\n    + (61) start: rc=0 (ok)\n   vip: migration-threshold=1000000\n    + (16) probe: rc=0 (ok)\n    + (21) stop: rc=0 (ok)\n   gfs: migration-threshold=1000000\n    + (43) probe: rc=0 (ok)\n    + (46) stop: rc=0 (ok)\n   net: migration-threshold=1000000\n    + (65) probe: rc=0 (ok)\n    + (66) stop: rc=0 (ok)\n   color: migration-threshold=1000000\n    + (70) probe: rc=0 (ok)\n    + (71) stop: rc=0 (ok)"
}
```

#### 2.11.2 获取指令列表

说明：获取预定义的指令列表

URI：/api/v1/haclusters/1/commands

Method：GET

请求参数：无

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| data                          | object                             | 预定义的指令集，map格式               |

响应示例:

```
{
    "action": True,
    "data": {
        "1": "crm_mon -1 -o",
        "2": "crm_simulate -Ls",
        "3": "pcs config show",
        "4": "corosync-cfgtool -s",
        "5": "crm configure verify"
    }
}
```


### 2.12 本机操作

#### 2.12.1 本机 HA 操作

说明：本机HA操作

URI：/api/v1/haclusters/1/localnodes/:action

Method：PUT

请求参数：

| 参数名称           | 是否必填          | 传入方式            | 参数类型           | 参数说明                       |
| :----------------- | ---------------- | :----------------- | :----------------- | ------------------------------ |
| action             | 是               | path               | string             | 操作类型：start、stop、restart  |

响应格式：json

响应参数：

| 参数名称                      | 参数类型                            | 参数说明                             |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | 返回结果状态                          |
| info                          | string                             | 执行结果提示信息                      |

响应示例:

```
{
    "action": True,
    info": "Action on node success"
}
```

### 2.13 集群管理

#### 2.13.1 创建集群

URI：/api/v1/managec/cluster_setup

Method：POST

##### 请求参数

| 参数名称        | 位置 | 类型     | 是否必选 | 说明                  |
| --------------- | ---- | -------- | -------- | --------------------- |
| body            | body | object   | 否       | none                  |
| >>cluster_name  | body | string   | 是       | 集群名称              |
| >>data          | body | [object] | 是       | 创建集群的节点数据    |
| >>>>nodeid      | body | integer  | 是       | 节点id                |
| >>>>name        | body | string   | 是       | 节点名称              |
| >>>>password    | body | string   | 是       | 节点hacluster用户密码 |
| >>>>ring0_adddr | body | [object] | 是       | 节点心跳ip            |

```
{
  "cluster_name": "hacluster",
  "data": [
    {
      "nodeid": 1,
      "name": "sp3-2",
      "password": "kylinha10!)",
      "ring_addr":[
          {
              "ring":"ring0_addr",
              "ip":"192.168.174.133"
          },
          {
              "ring":"ring1_addr",
              "ip":"172.30.30.93"
          }
      ]
    },
    {
      "nodeid": 2,
      "name": "sp3-4",
      "password": "kylinha10!)",
      "ring_addr":[
          {
              "ring":"ring0_addr",
              "ip":"192.168.174.133"
          },
          {
              "ring":"ring1_addr",
              "ip":"172.30.30.93"
          }
      ]
    }
  ]
}
```



##### 返回示例：成功

```
{
    "action": true,
    "message": "集群创建成功"
}
```

##### 返回结果

| 状态码 | 状态码含义 | 说明 | 数据模型 |
| ------ | ---------- | ---- | -------- |
| 200    | OK         | 成功 | Inline   |

##### 返回数据结构

状态码 **200**

| 参数名称  | 类型    | 是否必选 | 约束 | 说明                 |
| --------- | ------- | -------- | ---- | -------------------- |
| >>action  | boolean | true     | none | 创建集群操作是否成功 |
| >>message | string  | true     | none | 提示信息             |



#### 2.13.2 摧毁集群

URI：/api/v1/managec/cluster_destroy

Method：POST

##### Body请求参数

```
{
  "cluster_name": [
    "hacluster"
  ]
}
```

| 参数名称       | 位置 | 类型   | 是否必选 | 说明     |
| -------------- | ---- | ------ | -------- | -------- |
| body           | body | object | 否       |          |
| >>cluster_name | body | string | 是       | 集群名称 |

##### 返回示例：成功

```
{
  "action": true,
  "data": [
    true
  ],
  "clusters": [],
  "detailInfo": []
}
```

##### 返回结果

| 状态码 | 状态码含义 | 说明 | 数据模型 |
| ------ | ---------- | ---- | -------- |
| 200    | OK         | 成功 | Inline   |

##### 返回数据结构

状态码**200**

| 名称         | 类型      | 必选 | 约束 | 说明                       |
| ------------ | --------- | ---- | ---- | -------------------------- |
| >>action     | boolean   | true | none | 摧毁操作是否成功           |
| >>data       | [boolean] | true | none | 每个集群摧毁操作的结果汇总 |
| >>cluster    | [string]  | true | none | 摧毁失败的集群列表         |
| >>detailInfo | [string]  | true | none | 摧毁失败集群的详细信息     |

#### 2.13.3 移除集群

URL: /api/v1/managec/cluster_remove

Method：POST

##### Body请求参数

```
{
  "cluster_name": [
    "hacluster"
  ]
}
```

| 名称           | 位置 | 类型     | 必选 | 说明     |
| -------------- | ---- | -------- | ---- | -------- |
| body           | body | object   | 否   | none     |
| >>cluster_name | body | [string] | 是   | 集群名称 |

##### 返回示例

```
{
  "action": true,
  "faild_cluster": [],
  "data": [
    true
  ]
}
```

##### 返回结果

| 状态码 | 状态码含义 | 说明 | 数据模型 |
| ------ | ---------- | ---- | -------- |
| 200    | OK         | 成功 | Inline   |

##### 返回数据结构

状态码 **200**

| 名称            | 类型      | 必选 | 约束 | 说明               |
| --------------- | --------- | ---- | ---- | ------------------ |
| >>action        | boolean   | true | none | 移除操作是否成功   |
| >>faild cluster | [string]  | true | none | 移除失败的集群名称 |
| >>data          | [boolean] | true | none | 移除集群结果汇总   |

### 2.14 节点管理

### 2.14.1 添加节点

URL: /api/v1/managec/add_nodes

Method：POST

##### Body请求参数

```
{
  "cluster_name": "hacluster",
  "data": [
    {
      "name": "sp3-3",
      "password": "kylinha10!)",
      "ring_addr":[
          {
              "ring":"ring0_addr",
              "ip":"192.168.174.133"
          },
          {
              "ring":"ring1_addr",
              "ip":"172.30.30.93"
          }
      ]
    }
  ]
}
```

##### 请求参数

| 名称            | 位置 | 类型     | 必选 | 说明                  |
| --------------- | ---- | -------- | ---- | --------------------- |
| body            | body | object   | 否   | none                  |
| >> cluster_name | body | string   | 是   | 集群名称              |
| >>data          | body | [object] | 是   | none                  |
| >>name          | body | string   | 否   | 节点名称              |
| >>password      | body | string   | 否   | 节点hacluster用户密码 |
| >>ring_addr     | body | [object] | 否   | 节点心跳ip            |

##### 返回示例

```
成功
{
  "action": true,
  "message": "添加节点成功"
}
失败
{
  "action": false,
  "error": "添加节点失败",
  "detailInfo": "Error: Node name 'sp3-3' is already used by existing nodes; please, use other name\nError: Node address '172.30.230.105' is already used by existing nodes; please, use other address\nError: sp3-3: Running cluster services: 'corosync', 'pacemaker', the host seems to be in a cluster already, use --force to override\nError: sp3-3: Cluster configuration files found, the host seems to be in a cluster already, use --force to override\nError: Some nodes are already in a cluster. Enforcing this will destroy existing cluster on those nodes. You should remove the nodes from their clusters instead to keep the clusters working properly, use --force to override\nError: Errors have occurred, therefore pcs is unable to continue"
}
```

##### 返回结果

| 状态码 | 状态码含义 | 说明 | 数据模型 |
| ------ | ---------- | ---- | -------- |
| 200    | OK         | 成功 | Inline   |

##### 返回数据结构

状态码 **200**

| 名称      | 类型    | 必选 | 约束 | 说明                 |
| --------- | ------- | ---- | ---- | -------------------- |
| >>action  | boolean | true | none | 添加节点操作是否成功 |
| >>message | string  | true | none | 提示信息             |

