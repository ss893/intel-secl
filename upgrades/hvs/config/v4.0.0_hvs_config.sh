#!/bin/bash

SERVICE_NAME=hvs
CONFIG_FILE="/etc/$SERVICE_NAME/config.yml"
LOG_PATH=/var/log/$SERVICE_NAME

echo "Starting $SERVICE_NAME config upgrade to v4.0.0"
# Add ENABLE_EKCERT_REVOKE_CHECK setting to config.yml
grep -q 'enable-ekcert-revoke-check' $CONFIG_FILE || echo 'enable-ekcert-revoke-check: false' >>$CONFIG_FILE

echo "Completed $SERVICE_NAME config upgrade to v4.0.0"
