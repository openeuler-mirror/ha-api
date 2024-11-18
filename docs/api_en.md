# Rest API Document

## 1 API Description

### 1.1 API Brief

This document describes HA manage platform backend REST API interface.

### 1.2 Status Code

| Status Code | description                           |
| :---------: | :------------------------------------ |
|  2xx        | process normal                        |
|  3xx        | redirect to new URL                   |
|  4xx        | client request error                  |
|  5xx        | server process error                  |

### 1.3 Error Response

All error response use the unified return format as follow:

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| error                         | string                             | result error information             |

Error Response Example:

```
{
    "action":false,
    "error":"error info"
}
```


### 1.4 Update Records

| Index | Update Details      |    date      |
| :---: | :------------------ | :----------: |
|  1    | Initial version     | 2021.01.11   |


<div STYLE="page-break-after: always;"></div>



## 2 API Interface Definicaton

### 2.1 Cluster

#### 2.1.1 Get Cluster Attribution

Description: Get Cluster Attribution

URI：/api/v1/haclusters/1

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                    | Parameter Type                        | Parameter Description               |
| :-------------------------------- | :------------------------------------ | ----------------------------------- |
| action                            | bool                                  | result status                       |
| data                              | object                                | normal result data                  |
| error                             | string                                | error result information            |
| parameters                        | object                                | parameter details of every module   |
| shortdesc                         | string                                | short description                   |
| version                           | string                                | version                             |
| nodecount                         | int                                   | node count                          |
| isconfig                          | bool                                  | is config                           |
| longdesc                          | string                                | full description                    |


Response Example:

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



#### 2.1.2 Modify Cluster Attribution

Description: Modify Cluster Attribution

URI：/api/v1/haclusters/1

Method：PUT

Request Parameters:

| Parameter Name             | Must Provide     | Parameter Location | Parameter Type     | Parameter Description                   |
| :------------------------- | ---------------- | :----------------- | :----------------- | --------------------------------------- |
| no-quorum-policy(example)  | yes              | json               | string             | cluster attribute name and new config   |

Request Example:

```
{
    "no-quorum-policy": "stop"
}
```

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":true,
    "info":"Save crm metadata success"
}
```

#### 2.1.3 Cluster Operation

Description: cluster operation

URI：/api/v1/haclusters/1/:action

Method：PUT

Request Parameters:

| Parameter Name    | Must Provide    | Parameter Location  | Parameter Type  | Parameter Description                       |
| :---------------- | --------------- | :------------------ | :-------------- | ------------------------------------------- |
| action            | yes             | path                | string          | cluster operation. start, stop or restart   |
| nodeid            | yes             | json                | string          | node id                                     |
| nodeip            | yes             | json                | string          | node ip                                     |
| password          | yes             | json                | string          | node password                               |

Request Example:

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

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":true,
    "info":"Save crm metadata success"`
}
```

### 2.2 Node

#### 2.2.1 Get Node List

Description: get node list

URI：/api/haclusters/1/nodes

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | string                             | node information list                |

Response Example:

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

#### 2.2.2 Get Single Node Information

Description: get single node information

URI：/api/haclusters/1/nodes/:nodeid

Method：GET

Request Parameters:

| Parameter Name    | Must Provide    | Parameter Location | Parameter Type       | Parameter Description   |
| :---------------- | --------------- | :----------------- | :------------------- | ----------------------- |
| nodeid            | yes             | path               | string               | node id                 |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | string                             | resutl data                          |
| ips                           | array                              | node ip list                         |

Response Example:

```
{   "action":true,
    "data": {
        'ips': ['10.1.110.188', '192.168.100.188']
    }
}
```


#### 2.2.3 Node Operation

Description: Node Operation

URI：/api/v1/haclusters/1/nodes/:node_id/:action

Method：PUT

Request Parameters:

| Parameter Name    | Must Provide    | Parameter Location | Parameter Type  | Parameter Description             |
| :---------------- | --------------- | :----------------- | :-------------- | --------------------------------- |
| node_id           | yes             | path               | string          | node id                           |
| action            | yes             | path               | string          | node operation, including unstandby, standby, stop, start and restart     |

Request Example:

```
// start, stop and restart needs user password
{
    "password": "12345678",
}
```

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{   "action":true,
    "info":"Change node status success"
}
```

### 2.3 Resource

#### 2.3.1 Get Resource List

Description: get resource list

URI：/api/v1/haclusters/1/resources

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | result data                          |

Response Example:

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


#### 2.3.2 Add Resource

Description: add resource

URI：/api/v1/haclusters/1/resources

Method：POST

Request Parameters:

| Parameter Name   | Must Provide    | Parameter Location  | Parameter Type  | Parameter Description      |
| :--------------- | --------------- | :------------------ | :-------------- | -------------------------- |
| category         | yes             | json                | string          | resource category          |

Request Example:

```
// primitive resource
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
// clone resource
{
  "category": "clone",    
    "id":"test5",
    "rsc_id":"test4",
    "meta_attributes":{
        "target-role":"Stopped"
    }
}
// group resource
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

Response Format: json

Response Parameters:

| Parameter Name            | Parameter Type                | Parameter Description           |
| :------------------------ | :---------------------------- | ------------------------------- |
| action                    | bool                          | result status                   |
| info                      | string                        | result information              |

Response Example:

```
{
    "action":True,
    'info':"Add primitive/clone/group resource success"
}
```

#### 2.3.3 Single Resource Operation

Description: single resource operation

URI：/api/v1/haclusters/1/resources/:rsc_id/:action

Method：PUT

Request Parameters:

| Parameter Name    | Must Provide    | Parameter Location  | Parameter Type   | Parameter Description             |
| :---------------- | --------------- | :------------------ | :--------------- | --------------------------------- |
| rsc_id            | yes             | path                | string           | resource id                       |
| action            | yes             | path                | string           | resource operation, including start, stop, delete, cleanup, migrate, unmigrate, location, order, colocation     |

Request Example:

```
// start, stop, delete and cleanup
{}
// migrate
{
    "is_force": True,
    "to_node": "ns187",
    "period": "PYMDTHM3S"
}
// unmigrate
{
    "rsc_id": "kk1",
    "is_all_rscs":False
}
// location
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
// colocation
{
    "same_node": ["test1234"],
    "diff_node": ["group_tomcat"]
}
// order
{
    "before_rscs": ["test1234"],
    "after_rscs": ["group-fs-ps"]
}
```

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description              |
| :---------------------------- | :--------------------------------- | ---------------------------------- |
| action                        | bool                               | result status                      |
| info                          | string                             | result information                 |
| error                         | string                             | operation error                    |

Response Example:

```
{   "action":true,
    "info":"Action on resource success"
}
```


#### 2.3.4 Get All Resource Creation Data

Description: get all resource creation data

URI：/api/v1/haclusters/1/metas/:rsc_class/:rsc_type/:rsc_provider

Method：GET

Request Parameters:

| Parameter Name    | Must Provide    | Parameter Location  | Parameter Type | Parameter Description             |
| :---------------- | --------------- | :------------------ | :------------- | --------------------------------- |
| rsc_class         | yes             | path                | string         | resource class                    |
| rsc_type          | yes             | path                | string         | resource type                     |
| rsc_provider      | yes             | path                | string         | resource provider                 |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | resource creation data               |

Response Example:

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

#### 2.3.5 Get All Resource Type

Description: get all tesource type

URI：/api/v1/haclusters/1/metas

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | all resource type                    |

Response Example:

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

#### 2.3.6 Get Resource Meta Attributes

Description: get resource meta attributes

URI：/api/v1/haclusters/1/resources/meta_attributes/:catagory

Method：GET

Request Parameters:

| Parameter Name   | Must Provide   | Parameter Location  | Parameter Type | Parameter Description                        |
| :--------------- | -------------- | :------------------ | :------------- | -------------------------------------------- |
| catagory         | yes            | path                | string         | resource type, clone, primitive or group     |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | resource meta attributes             |

Response Example:

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

#### 2.3.7 Get Resource Relation Information

Description: get resource relation information

URI：/api/v1/haclusters/1/resources/:rsc_id/relations/:relation

Method：GET

Request Parameters:

| Parameter Name   | Must Provide   | Parameter Location | Parameter Type  | Parameter Description                          |
| :--------------- | -------------- | :----------------- | :-------------- | ---------------------------------------------- |
| rsc_id           | yes            | path               | string          | resource id                                    |
| relation         | yes            | path               | string          | relation type, location, order or colocation   |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | relation config                      |

Response Example:

```
// location
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

// colocation
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
// order
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

#### 2.3.8 Modify Resource

Description: modify resource

URI：/api/v1/haclusters/1/resources/:rsc_id

Method：PUT

Request Parameters:

| Parameter Name   | Must Provide   | Parameter Location | Parameter Type  | Parameter Description      |
| :--------------- | -------------- | :----------------- | :-------------- | -------------------------- |
| rsc_id           | yes            | path               | string          | resource id                |

Request Example:

```
// primitive
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
// group
{
    "id":"group1",
    "category":"group",
    "rscs":["iscisi", "test1" ],
    "meta_attributes":{
        "target-role":"Started"
    }
}
// clone
{
    "id":"clone1",
    "category":"clone",
    "rsc_id":"ip1",
    "meta_attributes":{
        "target-role":"Stopped"
    }
}
```


Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | object                             | result information                   |

Response Example:

```
// primitive
{
    'action':True,
    'info':"Set rsc info success"
}
// group
{
    'action':True,
    'info':"Set rsc info success"
}
// clone
{
    'action':True,
    'info':"Set rsc info success"
}
```

#### 2.3.9 Get Reource Information

Description: get reource information

URI：/api/v1/haclusters/1/resources/:rsc_id

Method：GET

Request Parameters:

| Parameter Name   | Must Provide   | Parameter Location  | Parameter Type  | Parameter Description      |
| :--------------- | -------------- | :------------------ | :-------------- | -------------------------- |
| rsc_id           | yes            | path                | string          | resource id                |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | resource information                 |

Response Example:

```
// primitive
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
// group
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
// clone
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

### 2.4 login

#### 2.4.1 login

Description: user login

URI：/api/v1/login

Method：POST

Request Parameters:

| Parameter Name   | Must Provide   | Parameter Location | Parameter Type  | Parameter Description         |
| :--------------- | -------------- | :----------------- | :-------------- | ----------------------------- |
| username         | yes            | json               | string          | user name                     |
| password         | yes            | json               | string          | password                      |

Request Example:

```
{
    "username": "hacluster",
    "password": "12345678"
}
```

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |

Response Example:

```
{
    "action": true
}
```

### 2.5 Alarm

#### 2.5.1 Get Alarm Information

Description: get alarm information

URI：/api/v1/haclusters/1/alarms

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | int                                | result data                          |
| sender                        | string                             | sender email                         |
| smtp                          | string                             | email smtp                           |
| flag                          | bool                               | switch flag                          |
| receiver                      | array                              | receiver email list                  |
| password                      | string                             | email password                       |
| port                          | string                             | email port                           |

Response Example:

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

#### 2.5.1 Create Alarm

Description: create alarm

URI：/api/v1/haclusters/1/alarms

Method：POST

Request Parameters:

| Parameter Name   | Must Provide   | Parameter Location | Parameter Type  | Parameter Description      |
| :--------------- | -------------- | :----------------- | :-------------- | -------------------------- |
| sender           | yes            | json               | string          | sender email               |
| smtp             | yes            | json               | string          | email smtp                 |
| flag             | yes            | json               | bool            | switch flag                |
| receiver         | yes            | json               | array           | reciever email list        |
| password         | yes            | json               | string          | email password             |
| port             | yes            | json               | string          | email port                 |

Request Example:

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

Response Format: json

Response Parameters:

| Parameter Name            | Parameter Type                | Parameter Description           |
| :------------------------ | :---------------------------- | ------------------------------- |
| action                    | bool                          | result status                   |
| info                      | string                        | result information              |

Response Example:

```
{
    "action":True,
    'info':"Set alarm success"
}
```

#### 2.5.1 Delete Alarm

TODO:

### 2.6 Heartbeat

#### 2.6.1 Get Network Heartbeat Information

Description: get network heartbeat information

URI：/api/v1/haclusters/1/configs

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                         |
| :---------------------------- | :--------------------------------- | --------------------------------------------- |
| action                        | bool                               | result status                                 |
| data                          | int                                | heartbeat status, 0: normal, not 0: abnormal  |
| hbaddrs1                      | object                             | single heartbeat config                       |
| hbaddrs2                      | object                             | redundance heartbeat config                   |
| ip                            | string                             | node ip                                       |
| nodeid                        | string                             | node id                                       |
| hbaddrs2_enabled              | int                                | redundance heartbeat enabled flag             |

Response Example:

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

#### 2.6.2 Create/Modify Heartbeat Configuration

Description: create/modify heartbeat configuration

URI：/api/v1/haclusters/1/configs

Method：POST

Request Parameters:

| Parameter Name       | Must Provide  | Parameter Location | Parameter Type   | Parameter Description               |
| :------------------- | ------------- | :----------------- | :--------------- | ----------------------------------- |
| hbaddrs1             | yes           | json               | object           | single heartbeat config             |
| hbaddrs2             | no            | json               | object           | redundance heartbeat config         |
| ip                   | yes           | json               | string           | node ip                             |
| nodeid               | yes           | json               | string           | node id                             |
| hbaddrs2_enabled     | yes           | json               | int              | redundance heartbeat enabled flag   |

Request Example:

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

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":True,
    'info':"update success"
}
```

#### 2.6.3 Get Available Dist Heartbeat Device List

TODO: 

#### 2.6.4 Get Heartbeat Status

Description: get heartbeat status

URI：/api/v1/haclusters/1/hbstatus

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                         |
| :---------------------------- | :--------------------------------- | --------------------------------------------- |
| action                        | bool                               | result status                                 |
| data                          | int                                | heartbeat status, 0: normal, not 0: abnormal  |

Response Example:

```
{
    "action":True,
    "data": 0
}
```


### 2.7 DRBD

#### 2.7.1 Get DRBD Information

Description: get DRBD information

URI：/api/v1/haclusters/1/drbd

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                         |
| :---------------------------- | :--------------------------------- | --------------------------------------------- |
| action                        | bool                               | result status                                 |
| data                          | object                             | result information                            |
| dir                           | string                             | dir                                           |
| main_node_ip                  | string                             | main node ip                                  |
| main_node_name                | string                             | main node name                                |
| main_node_device              | string                             | main node device                              |
| slave_node_ip                 | string                             | slave node ip                                 |
| slave_node_name               | string                             | slave node name                               |
| slave_node_device             | string                             | slave node device                             |
| drbd_type                     | string                             | drbd type，0: master-slave, 1: master-master  |

Response Example:

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

#### 2.7.2 Create/Modify DRBD

Description: create/modify DRBD

URI：/api/v1/haclusters/1/drbd

Method：POST

Request Parameters:

| Parameter Name       | Must Provide   | Parameter Location | Parameter Type | Parameter Description                          |
| :------------------- | -------------- | :----------------- | :------------- | ---------------------------------------------- |
| dir                  | yes            | json               | string         | path                                           |
| main_node_ip         | yes            | json               | string         | main node ip                                   |
| main_node_name       | yes            | json               | string         | main node name                                 |
| main_node_device     | yes            | json               | string         | main node device                               |
| slave_node_ip        | yes            | json               | string         | slave node ip                                  |
| slave_node_name      | yes            | json               | string         | slave node name                                |
| slave_node_device    | yes            | json               | string         | slave node device                              |
| drbd_type            | yes            | json               | string         | drbd type，0: master-slave, 1: master-master   |

Request Example:

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

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":True,
    'info':"Save drbd configuration success"
}
```

#### 2.7.3 Get DRBD Status

TODO: 

#### 2.7.4 Delete DRBD

Description: Delete DRBD

URI：/api/v1/haclusters/1/drbd

Method：DELETE

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":True,
    'info':"Delete drbd configuration success"
}
```

### 2.8 GFS

#### 2.8.1 Get GFS Information

Description: get GFS information

URI：/api/v1/haclusters/1/gfs

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | result information                   |
| device                        | string                             | device path                          |
| max_num                       | int                                | max number                           |
| mount_dir                     | string                             | device mount directory               |

Response Example:

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


#### 2.8.2 Create/Modify GFS

Description: create/modify GFS

URI：/api/v1/haclusters/1/gfs

Method：POST

Request Parameters:

| Parameter Name     | Must Provide     | Parameter Location  | Parameter Type     | Parameter Description                  |
| :----------------- | ---------------- | :------------------ | :----------------- | ------------------------ |
| device             | yes              | json                | string             | device path              |
| max_num            | yes              | json                | int                | max number               |
| mount_dir          | yes              | json                | string             | device mount directory   |

Request Example:

```
{
    "device": "/dev/vdb",
    "max_num": 2,
    "mount_dir": "/drbd"
}
```

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":True,
    'info':"Save gfs configuration success"
}
```

#### 2.8.3 Delete GFS

Description: delete GFS

URI：/api/v1/haclusters/1/gfs

Method：DELETE

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action":True,
    'info':"Delete gfs success"
}
```

### 2.9 Script

#### 2.9.1 Generate Script

Description: generate resource control script

URI：/api/v1/haclusters/1/scripts

Method：POST

Request Parameters:

| Parameter Name     | Must Provide      | Parameter Location | Parameter Type     | Parameter Description                  |
| :----------------- | ----------------- | :----------------- | :----------------- | ------------------------ |
| name               | yes               | json               | string             | script name              |
| start              | yes               | json               | string             | start control script     |
| stop               | yes               | json               | string             | stop control script      |
| monitor            | yes               | json               | string             | monitor control script   |

Request Example:

```
{
    "name": "tomcat3",
    "start":"/usr/local/bin/tomcat start",
    "stop": "/usr/local/bin/tomcat stop",
    "monitor": "/usr/local/bin/tomcat monitor"
}
```

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action": true,
    "info": "Create script success!"
}
```


### 2.10 log

#### 2.10.1 Gererage Log

Description: Gererage cluster log and return log path

URI：/api/v1/haclusters/1/logs

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | result data                          |
| filepath                      | string                             | compressed log file path             |

Response Example:

```
{
    'action': True,
    'data': {
        "filepath": "static/neokylinha-log.tar"
    }
}
```

### 2.11 Commands

#### 2.11.1 Get Commands Result

Description: run pre-defined commands and get result

URI：/api/v1/haclusters/1/commands/:cmd_type

Method：GET

Request Parameters:

| Parameter Name     | Must Provide      | Parameter Location | Parameter Type     | Parameter Description        |
| :----------------- | ----------------- | :----------------- | :----------------- | ---------------------------- |
| cmd_type           | yes               | path               | int                | pre-defined command id       |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | string                             | command result                       |

Response Example:

```
{
    "action": true,
    "data": "Stack: corosync\nCurrent DC: ns175 (version 1.1.16-12.el7.1-94ff4df) - partition with quorum\nLast updated: Thu Jun  7 15:55:05 2018\nLast change: Wed Jun  6 10:56:10 2018 by root via mgmtd on ns175\n\n2 nodes configured\n19 resources configured (23 DISABLED)\n\nOnline: [ ns175 ]\nOFFLINE: [ ns176 ]\n\nActive resources:\n\n dummy\t(ocf::heartbeat:Dummy):\tStarted ns175\n test2\t(ocf::heartbeat:Dummy):\tStarted ns175\n\nOperations:\n* Node ns175:\n   dummy: migration-threshold=1000000\n    + (60) start: rc=0 (ok)\n   test2: migration-threshold=1000000\n    + (61) start: rc=0 (ok)\n   vip: migration-threshold=1000000\n    + (16) probe: rc=0 (ok)\n    + (21) stop: rc=0 (ok)\n   gfs: migration-threshold=1000000\n    + (43) probe: rc=0 (ok)\n    + (46) stop: rc=0 (ok)\n   net: migration-threshold=1000000\n    + (65) probe: rc=0 (ok)\n    + (66) stop: rc=0 (ok)\n   color: migration-threshold=1000000\n    + (70) probe: rc=0 (ok)\n    + (71) stop: rc=0 (ok)"
}
```

#### 2.11.2 Get Commands List

Description: get commands list

URI：/api/v1/haclusters/1/commands

Method：GET

Request Parameters: None

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| data                          | object                             | pre-defined commands list            |

Response Example:

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


### 2.12 Local HA Operation

#### 2.12.1 Local HA Operation

Description: local HA operation

URI：/api/v1/haclusters/1/localnodes/:action

Method：PUT

Request Parameters:

| Parameter Name     | Must Provide     | Parameter Location | Parameter Type     | Parameter Description                   |
| :----------------- | ---------------- | :----------------- | :----------------- | --------------------------------------- |
| action             | yes              | path               | string             | operation type, start, stop or restart  |

Response Format: json

Response Parameters:

| Parameter Name                | Parameter Type                     | Parameter Description                |
| :---------------------------- | :--------------------------------- | ------------------------------------ |
| action                        | bool                               | result status                        |
| info                          | string                             | result information                   |

Response Example:

```
{
    "action": True,
    info": "Action on node success"
}
```
