#!/bin/bash

SERVICE_USERNAME=ihub
COMPONENT_NAME=$SERVICE_USERNAME
PRODUCT_HOME=/opt/$COMPONENT_NAME
INSTANCE_NAME=${INSTANCE_NAME:-$COMPONENT_NAME}
BACKUP_PATH=${BACKUP_PATH:-"/tmp/"}
BACKUP_DIR=${BACKUP_PATH}${SERVICE_USERNAME}_backup
LOG_PATH=/var/log/$COMPONENT_NAME

echo "Starting $COMPONENT_NAME config upgrade to v3.6.0"
# Update config file
echo "Using HVS Base Url $HVS_BASE_URL"
echo "Using SHVS Base Url $SHVS_BASE_URL"
./$COMPONENT_NAME setup attestation-service-connection
if [ $? -ne 0 ]; then
  echo "Failed to update config to v3.6.0"
  exit 1
fi

# Install systemd script
SERVICE_FILE=$SERVICE_USERNAME@.service
cp $SERVICE_USERNAME.service $PRODUCT_HOME/$SERVICE_FILE && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME/$SERVICE_FILE && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME

chmod 640 $LOG_PATH/*
chmod 740 $LOG_PATH

# Enable systemd service
systemctl disable $SERVICE_USERNAME.service >/dev/null 2>&1
systemctl disable $PRODUCT_HOME/$SERVICE_FILE >/dev/null 2>&1
systemctl enable $PRODUCT_HOME/$SERVICE_FILE
systemctl disable $COMPONENT_NAME@$INSTANCE_NAME >/dev/null 2>&1
systemctl enable $COMPONENT_NAME@$INSTANCE_NAME
systemctl daemon-reload

# Backup service file post service stop
mv $PRODUCT_HOME/ihub.service "$BACKUP_DIR"/

echo "Completed $COMPONENT_NAME config upgrade to v3.6.0"
