#!/bin/bash

# test command
go test -v

if [ $? -eq 0 ]
then
  echo "<<<Test PASS>>>"
  touch /var/tmp/success
  exit 0
else
  echo "<<<TEST FAIL>>>" >&2
  exit 1
fi