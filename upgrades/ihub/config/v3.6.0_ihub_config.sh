#!/bin/bash

SERVICE_USERNAME=ihub
COMPONENT_NAME=$SERVICE_USERNAME
PRODUCT_HOME=/opt/$COMPONENT_NAME
INSTANCE_NAME=${INSTANCE_NAME:-$COMPONENT_NAME}

echo "Starting $COMPONENT_NAME config upgrade to v3.6.0"
# Install systemd script
SERVICE_FILE=$SERVICE_USERNAME@.service
cp $SERVICE_USERNAME.service $PRODUCT_HOME/$SERVICE_FILE && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME/$SERVICE_FILE && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME

# Enable systemd service
systemctl disable $PRODUCT_HOME/$SERVICE_FILE >/dev/null 2>&1
systemctl enable $PRODUCT_HOME/$SERVICE_FILE
systemctl disable $COMPONENT_NAME@$INSTANCE_NAME >/dev/null 2>&1
systemctl enable $COMPONENT_NAME@$INSTANCE_NAME
systemctl daemon-reload

echo "Completed $COMPONENT_NAME config upgrade to v3.6.0"
