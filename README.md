# psmb
[![Build Status](https://travis-ci.org/taka-wang/psmb.svg?branch=dev)](https://travis-ci.org/taka-wang/psmb)
[![GoDoc](https://godoc.org/github.com/taka-wang/psmb?status.svg)](http://godoc.org/github.com/taka-wang/psmb)
[![GitHub tag](https://img.shields.io/github/tag/taka-wang/psmb.svg)](https://github.com/taka-wang/psmb/tags) 
[![Release](https://img.shields.io/github/release/taka-wang/psmb.svg)](https://github.com/taka-wang/psmb/releases/latest)


Proactive service for [modbusd](https://github.com/taka-wang/modbusd)

---

# Unit tests

- [binary](binary_test.go)
- [upstream](upstream_test.go)
- [downstream](downstream_test.go)

---

# Module
- [scheduler](https://github.com/taka-wang/gocron)
- request parser
- command builder
- post processing (event mapper)
    - filter
    - on changed
- zmq sub from upstream
- zmq sub from downstream
- logger
- database?

---
## Docker 

### Build images
```bash
docker build -t takawang/mbd test/mb/.

docker build -t takawang/psmb .

docker build -t takawang/psmb-goclient test/goclient/.

```

### Run images
```bash
```
