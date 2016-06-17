# psmb

Proactive service for [modbusd](https://github.com/taka-wang/modbusd)


# Request type

## Polling requests
- export
- load requests
- read all requests with status
- delete all requests
- enable all requests
- disable all requests
- add request
- update request
- delete request
- enable request
- disable request

## Rule requests
- export
- load rules
- read all rules with status
- delete all rules
- enable all rules
- disable all rules
- add rule
- update rule
- delete rule
- enable rule
- disable rule

## One-shot request
- trigger
    - read
    - write

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