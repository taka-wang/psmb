# psmb
[![Build Status](https://travis-ci.org/taka-wang/psmb.svg?branch=dev)](https://travis-ci.org/taka-wang/psmb)

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

# Note:
- if scheduler is stop, trigger requests directly.
- handle default port if not set
- handle data length
    - read: default = 1
    - write: check length with data, or set it automatically
- if the length of the response data equal to 1, should we put it data 'array'
- check timeout interval range
- check polling interval range
