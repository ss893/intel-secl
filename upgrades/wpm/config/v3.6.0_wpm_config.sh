#!/bin/bash

SERVICE_USERNAME=wpm
COMPONENT_NAME=workload-policy-manager
CONFIG_PATH=/etc/$COMPONENT_NAME

echo "Starting $COMPONENT_NAME config upgrade to v3.6.0"
# Update config file
../config-upgrade v3.6.0_config.tmpl $1/config.yml $CONFIG_PATH/config.yml
if [ $? -ne 0 ]; then
  echo "Failed to update config to v3.6.0"
  exit 1
fi

echo "Completed $COMPONENT_NAME config upgrade to v3.6.0"
