# psmb
[![Build Status](https://travis-ci.org/taka-wang/psmb.svg?branch=dev)](https://travis-ci.org/taka-wang/psmb)
[![GoDoc](https://godoc.org/github.com/taka-wang/psmb?status.svg)](http://godoc.org/github.com/taka-wang/psmb)

Proactive service for [modbusd](https://github.com/taka-wang/modbusd)

# Module
- scheduler
- request parser
- command builder
- post processing (event mapper)
    - filter
    - on changed
- zmq sub from upstream
- zmq sub from downstream
- logger
- database?
