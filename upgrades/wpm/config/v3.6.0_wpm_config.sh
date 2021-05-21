#!/bin/bash

SERVICE_USERNAME=wpm
COMPONENT_NAME=workload-policy-manager
CONFIG_PATH=/etc/$COMPONENT_NAME
LOG_PATH=/var/log/$COMPONENT_NAME

# New Directories
FLAVORS="/opt/workload-policy-manager/flavors"
VM_IMAGES_PATH="/opt/workload-policy-manager/vm-images"
ENCRYPTED_VM_IMAGES_PATH="/opt/workload-policy-manager/encrypted-vm-images"

echo "Starting $COMPONENT_NAME config upgrade to v3.6.0"
# Update config file
./config-upgrade ./config/v3.6.0_config.tmpl $1/config.yml $CONFIG_PATH/config.yml
if [ $? -ne 0 ]; then
  echo "Failed to update config to v3.6.0"
  exit 1
fi

chmod 640 $LOG_PATH/*
chmod 740 $LOG_PATH

FLAVOR_SIGNING_CERTS_PATH=$CONFIG_PATH/certs/flavorsign
echo "Renaming flavor-signing-cert.pem to flavor-signing.pem"
mv $FLAVOR_SIGNING_CERTS_PATH/flavor-signing-cert.pem $FLAVOR_SIGNING_CERTS_PATH/flavor-signing.pem
sed -i 's/flavor-signing-cert\.pem/flavor-signing\.pem/g' $CONFIG_PATH/config.yml
echo "Renaming flavor-signing-key.pem to flavor-signing.key"
mv $FLAVOR_SIGNING_CERTS_PATH/flavor-signing-key.pem $FLAVOR_SIGNING_CERTS_PATH/flavor-signing.key
sed -i 's/flavor-signing-key\.pem/flavor-signing\.key/g' $CONFIG_PATH/config.yml

# Update paths in the config file to the new path
sed -i 's/\/etc\/workload-policy-manager\//\/etc\/wpm\//g' $CONFIG_PATH/config.yml

echo "Creating new directories... "
for directory in $FLAVORS $VM_IMAGES_PATH $ENCRYPTED_VM_IMAGES_PATH; do
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
  chmod 700 $directory
done

echo "Completed $COMPONENT_NAME config upgrade to v3.6.0"
