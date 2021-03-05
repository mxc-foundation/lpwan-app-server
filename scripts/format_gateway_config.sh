#!/bin/bash

if [ "$1"x = ""x ]; then
  echo "Usage: ./format_gateway_config.sh CONFIG_FILE_NAME "
  exit
fi

PARENT_PATH=$(dirname $0)
CONFIG_FILE_PATH="$PARENT_PATH"/../static/gateway-config/"$1"

echo "Format config file $CONFIG_FILE_PATH "

cat $CONFIG_FILE_PATH | tr -d "\r\n\t\v\0" | sed 's/\"/\\"/g' >>$CONFIG_FILE_PATH.format
