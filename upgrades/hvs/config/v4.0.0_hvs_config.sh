#!/bin/bash

SERVICE_NAME=hvs
SERVICE_USERNAME=hvs
CONFIG_DIR=$2

echo "Starting $SERVICE_NAME config upgrade to v4.0.0"
TEMPLATES_PATH=$CONFIG_DIR/templates
SCHEMA_PATH=$CONFIG_DIR/schema

mkdir -p $TEMPLATES_PATH $SCHEMA_PATH

# Change permission only in case of container environment
if [ -f "/.container-env" ]; then
  # Copy Schema files
  cp /tmp/schema/*.json $SCHEMA_PATH/ && chmod 0600 $SCHEMA_PATH/*.json
  # Copy template files
  cp /tmp/templates/*.json $TEMPLATES_PATH/ && chmod 0600 $TEMPLATES_PATH/*.json
  echo "Completed $SERVICE_NAME config upgrade to v4.0.0"
  exit 0
fi

# Copy template files
cp -r templates/ $CONFIG_DIR/

# Copy Schema files
cp -r schema/ $CONFIG_DIR/

chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $SCHEMA_PATH/
chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $TEMPLATES_PATH/
echo "Completed $SERVICE_NAME config upgrade to v4.0.0"
