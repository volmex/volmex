# volmex

Volmex is a docker volume driver that allows to execute arbitrary commands before a volume is mounted (pre-mount hook).

## Context
Using a swarm cluster, we found that there is no satisfying solution to setup and use multi-host persistent volumes.

As there are numerous protocols and implementations for multi-host storage solutions, e.g. `rsync`, [syncthing](https://syncthing.net/), ... we decided to not implement a new protocol but strive for a more abstract solution.

Hence, we created a Docker volume plugin (v1) that simply executes a user defined command whenever a volume is mounted.

## General
Besides the 'pre-mount hook', volmex volumes are equivalent to named volumes in Docker (With a different storage base directory though - `/var/local/volmex`)


## Usage (Evaluation/Testing)

```
# get volmex
$ wget https://github.com/volmex/volmex/releases/download/v0.9/volmex-0.9.tar
$ tar xf volmex-0.9.1.tar

# create dir where local volumes are stored
$ sudo mkdir /var/local/volmex

# start volmex driver (later you can move this to a systemd service file)
$ sudo ./volmex-daemon

# create a volume 'foo' using the volmex driver and specify a command or shell script that shall be executed pre-mount
$ docker volume create \
  --driver volmex \
  --opt cmd="/usr/local/sbin/do-something" \
  foo
```

When the command is executed, the following variables are available in the command's environment:

+ `VOLMEX_NAME` = foo 
+ `VOLMEX_MOUNTPOINT` = /var/local/volmex/foo 
+ `VOLMEX_CMD` = /usr/local/sbin/do-something

## Install on systemd systems (Persistent)
+ Download and extract volmex

```
# cd /tmp
# wget https://github.com/volmex/volmex/releases/download/v0.9/volmex-0.9.tar
# tar xf volmex-0.9.1.tar

```

+ Install volmex

```
# mkdir -p /usr/lib/docker
# install -D -m 744 volmex-daemon /usr/lib/docker/volmex-daemon
# install -D -m 644 volmex.service /etc/systemd/system/volmex.service
# mkdir -p /var/local/volmex
# systemctl daemon-reload
```

+ Start/enable volmex

```
# systemctl enable volmex
# systemctl restart docker
```

+ Volmex should be started before Docker - which is configured in the service file. That's why it's sufficient to `systemctl restart docker`.

## Usage

- Create script that shall be executed on volume mount:

```
➜  ~ sudo vim /usr/local/sbin/do-something
...
#!/usr/bin/env bash
echo $VOLMEX_NAME
echo sleeping for 10 seconds...
sleep 10
echo done
...
➜  ~ sudo chmod +x /usr/local/sbin/do-something

```

- Create the volume using the volmex driver

```
➜  ~ docker volume create --driver volmex --opt cmd="/usr/local/sbin/do-something" foo
```

- Add content

```
➜  ~ sudo mkdir -p /var/local/volmex/foo
➜  ~ sudo vim /var/local/volmex/foo/volmex-test
...
volmex-test
...
```

- Create a container which uses the volume

```
➜  ~ docker run --rm -ti --volume foo:/foo alpine sh
/ # cat /foo/volmex-test
volmex-test
/ # exit
```

- Check log output of the executed command

```
➜  ~ journalctl -u volmex | tail
...
volmex-daemon[2115]: received 'Get' request for volume: foo
volmex-daemon[2115]: received 'Mount' request for volume: foo
volmex-daemon[2115]:         executing volmex command: /usr/local/sbin/do-something
volmex-daemon[2115]: foo
volmex-daemon[2115]: sleeping for 10 seconds...
volmex-daemon[2115]: done
volmex-daemon[2115]: received 'Unmount' request for volume: foo
...
```

