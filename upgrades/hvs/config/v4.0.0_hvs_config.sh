#!/bin/bash

SERVICE_NAME=hvs
SERVICE_USERNAME=hvs
CONFIG_DIR=$2

echo "Starting $SERVICE_NAME config upgrade to v4.0.0"
TEMPLATES_PATH=$CONFIG_DIR/templates
SCHEMA_PATH=$CONFIG_DIR/schema

mkdir -p $TEMPLATES_PATH $SCHEMA_PATH

# Copy template files
cp -r templates/ $CONFIG_DIR/ && chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $TEMPLATES_PATH/

# Copy Schema files
cp -r schema/ $CONFIG_DIR/ && chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $SCHEMA_PATH/

echo "Completed $SERVICE_NAME config upgrade to v4.0.0"
