#!/bin/bash

COMPONENT_NAME=kbs

if [ -f "/.container-env" ]; then
    set -a; source /etc/secret-volume/secrets.txt; set +a;
fi

echo "Starting $COMPONENT_NAME config upgrade to v4.0.0"
# Update config file
echo "Using KMIP Hostname $KMIP_HOSTNAME"
echo "Using KMIP Username $KMIP_USERNAME"
echo "Using KMIP Password $KMIP_PASSWORD"
./$COMPONENT_NAME setup update-service-config --force
if [ $? -ne 0 ]; then
  echo "Failed to update config to v4.0.0"
  exit 1
fi

echo "Completed $COMPONENT_NAME config upgrade to v4.0.0"
