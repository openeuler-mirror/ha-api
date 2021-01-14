# Architecture

This file describe HA-api management backend project architecture.


## Project Architecture

TODO: need a picture here

Briefly, HA-api wraps many HA manager commands such as pcs, crm_xxx, cibamdin, and ha-web project provides a easy-to-use web UI to monitor and control your HA cluster.

## Code Structure

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

controllers: REST api processer needed by Beego framework.
models: wraps HA cluster manager commands.
routers: bind URLs and controllers.
services: currently only session service.
settings: application settings.
utils: common utils.
views: static files.