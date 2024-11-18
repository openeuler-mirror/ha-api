# Build from Source

This document describes how to build HA-api management backend from code.

## Build

### Build Requirements

The following packages are needed to be installed on the host.

 - go >= 1.13
 - git
 - pam-devel
 - tar

### Build

Get the code first:

```
git clone https://gitee.com/openeuler/ha-api.git
cd ha-api/
git checkout -b release-v1.0 origin/release-v1.0
```

Simplely run `go build` to build the project.
```
go build
```
This will generate `ha-api`(or `ha-api.exe` on windows) executable file.

## Install and Usage

Before you use the HA-api backend, you need to install HA software first. Check [ha_install](./ha_install_en.md) for more information.

There is a install script under `tools` directory, it will build the executable file then install it and related script file. After installation, start service with systemd:

```
systemctl start ha-api
```
![run_ha-api](../pictures/run_ha-api.png)

Now you can use the REST API on default port 8080 to manage your HA cluster. You can use commands like `curl` to do that.

```
curl -s http://localhost:8080//api/v1/haclusters/1`
```
![hacluster_api_example](../pictures/hacluster_api_example.png)
