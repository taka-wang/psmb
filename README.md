# psmb

Proactive service for [modbusd](https://github.com/taka-wang/modbusd)


# Request type

## One-off request
- read coil/register (FC1,2,3,4)
- write coil/register(FC5,6,15,16)

## Polling requests
- import requests
- export requests
- read all requests with status
- delete all requests
- enable all requests
- disable all requests
- add request
- add requests
- update request
- delete request
- delete requests
- enable request
- disable request

## Filter requests
- import rules
- export rules
- read all rules with status
- delete all rules
- enable all rules
- disable all rules
- add rule
- update rule
- delete rule
- enable rule
- disable rule

*note*: if scheduler is stop, trigger request directly.

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


