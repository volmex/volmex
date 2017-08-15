# volmex

Volmex is a docker volume driver that allows to execute arbitrary commands before a volume is mounted (pre-mount hook).

[![Build Status](https://travis-ci.org/volmex/volmex.svg?branch=master)](https://travis-ci.org/volmex/volmex)

## Context
Using a swarm cluster, we found that there is no satisfying solution to setup and use multi-host persistent volumes.

As there are numerous protocols and implementations for multi-host storage solutions, e.g. `rsync`, [syncthing](https://syncthing.net/), ...
, we decided to not implement a new protocol but strive for a more abstract solution.

Hence, we created a Docker v1 volume plugin that simply executes a user defined command whenever a volume is mounted.

## Usage
```
# create dir where local volumes are synced to
$ sudo mkdir /var/local/volmex

# start volmex driver server (later you can move this to a systemd unit file)
$ sudo ./daemon

# create a volume using the volmex driver and specify a command or shell script
$ docker volume create \
  --driver volmex \
  --opt cmd="/usr/local/sbin/mount-someting" \
  foo

# check volume
$ docker volume ls
volmex      foo
```

When the command is executed, the following environment variables are available:

+ `VOLMEX_NAME` = foo 
+ `VOLMEX_MOUNTPOINT` = /var/local/volmex/foo 
+ `VOLMEX_CMD` = /usr/local/sbin/mount-something

## Install
TODO
