# psmb

[![Build Status](http://dds.cmwang.net/api/badges/taka-wang/psmb/status.svg)](http://dds.cmwang.net/taka-wang/psmb)
[![GoDoc](https://godoc.org/github.com/taka-wang/psmb?status.svg)](http://godoc.org/github.com/taka-wang/psmb)
[![GitHub tag](https://img.shields.io/github/tag/taka-wang/psmb.svg)](https://github.com/taka-wang/psmb/tags) 
[![Release](https://img.shields.io/github/release/taka-wang/psmb.svg)](https://github.com/taka-wang/psmb/releases/latest)

Proactive service for [modbusd](https://github.com/taka-wang/modbusd)

---

# Unit tests

- [binary](binary_test.go)
- [types](types_test.go)

---

# Contracts (Interfaces)

- ProactiveService: proactive service
- MbtcpReadTask: read/poll task map
- MbtcpSimpleTask: simple task map

---
## Docker 

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
# run modbus server
docker run -itd --name=slave takawang/c-modbus-slave
# run modbusd
docker run -v /tmp:/tmp --link slave -it --name=modbusd takawang/modbusd

# run psmb
docker run -v /tmp:/tmp -itd takawang/psmb
#docker run -v /tmp:/tmp -it takawang/psmb /bin/bash

# run dummy-srv
docker run -v /tmp:/tmp --link slave -it takawang/dummy-srv

```

### Deployment Diagram

![deployment](image/deployment.png)

