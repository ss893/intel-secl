#!/bin/bash

SERVICE_USERNAME=wpm
COMPONENT_NAME=workload-policy-manager
CONFIG_PATH=/etc/$COMPONENT_NAME

echo "Starting $COMPONENT_NAME config upgrade to v3.6.0"
# Update config file
./config-upgrade ./config/v3.6.0_config.tmpl $1/config.yml $CONFIG_PATH/config.yml
if [ $? -ne 0 ]; then
  echo "Failed to update config to v3.6.0"
  exit 1
fi

FLAVOR_SIGNING_CERTS_PATH=$CONFIG_PATH/certs/flavorsign
echo "Renaming flavor-signing-cert.pem to flavor-signing.pem"
mv $FLAVOR_SIGNING_CERTS_PATH/flavor-signing-cert.pem $FLAVOR_SIGNING_CERTS_PATH/flavor-signing.pem
echo "Renaming flavor-signing-key.pem to flavor-signing.key"
mv $FLAVOR_SIGNING_CERTS_PATH/flavor-signing-key.pem $FLAVOR_SIGNING_CERTS_PATH/flavor-signing.key

echo "Completed $COMPONENT_NAME config upgrade to v3.6.0"
