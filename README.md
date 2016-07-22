# psmb

[![Build Status](http://dev.cmwang.net/api/badges/taka-wang/psmb/status.svg)](http://dev.cmwang.net/taka-wang/psmb)
[![GoDoc](https://godoc.org/github.com/taka-wang/psmb?status.svg)](http://godoc.org/github.com/taka-wang/psmb)

Proactive service for [modbusd](https://github.com/taka-wang/modbusd)

---

## // Contracts (Interfaces)

- ProactiveService: proactive service
- MbtcpReadTask: read/poll task map
- MbtcpSimpleTask: simple task map


## // Docker 

### Docker Compose

```bash
docker-compose build  --pull
docker-compose up --abort-on-container-exit
```

### Build images manually

```bash
# build psmb image
docker build -t takawang/psmb .
```

### Run images

```bash

# run modbus slave
docker run -itd --name=slave takawang/c-modbus-slave

# run modbusd
docker run -v /tmp:/tmp --link slave -it --name=modbusd takawang/modbusd

# run psmb
docker run -v /tmp:/tmp -itd takawang/psmb

# run dummy-srv
docker run -v /tmp:/tmp --link slave -it takawang/dummy-srv
```

## // Continuous Integration

I do continuous integration and build docker images after git push by self-hosted drone.io server and [dockerhub]((https://hub.docker.com/r/takawang/c-modbus-slave/)) service.


## // Deployment Diagram

![deployment](image/deployment.png)

