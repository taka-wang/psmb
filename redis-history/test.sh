#!/bin/bash

#touch /root/test 2> /dev/null

go test -v

if [ $? -eq 0 ]
then
  echo "Successfully"
  touch /tmp/success
  exit 0
else
  echo "Fail" >&2
  exit 1
fi