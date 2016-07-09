#!/bin/bash
/usr/bin/modbusd /etc/modbusd/modbusd.json &
/usr/bin/server &
/usr/bin/server 503

