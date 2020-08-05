#!/usr/bin/env bash

set -e

if [ $# -eq 0 ]
then
  /app/http
fi

exec $@
