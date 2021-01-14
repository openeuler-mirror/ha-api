# Build from Source

This document describes how to build HA-api managment backend from code.

## Build

### Build Requirements

Only golang is needed to be installed on the host.

 - go >= 1.13

To get code from gitee, you alse neet to install `git`.

```
yum install git
```

### Build

Get the code first:

```
git clone https://gitee.com/openeuler/ha-api.git
cd ha-api/
git checkout -b go-api origin/go-api
```

Simplely run `go build` to build the project.
```
go build
```
This will generate `ha-api`(or `ha-api.exe` on windows) executable file.

## Usage

Before you use the HA-api backend, you need to install HA software first. Check [ha_install](./ha_install_en.md) for more information.

Run the server by `./ha-api` after build:ha_
```
./ha-api
```
![run_ha-api](../pictures/run_ha-api.png)

Now you can use the REST API on default port 8080 to manage your HA cluster. You can use commands like `curl` to do that.

```
curl -s http://localhost:8080//api/v1/haclusters/1`
```
![hacluster_api_example](../pictures/hacluster_api_example.png)