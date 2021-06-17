#!/bin/bash

SERVICE_NAME=hvs
SERVICE_USERNAME=hvs
CONFIG_DIR=$2

echo "Starting $SERVICE_NAME config upgrade to v4.0.0"
TEMPLATES_PATH=$CONFIG_DIR/templates
SCHEMA_PATH=$CONFIG_DIR/schema

mkdir -p $TEMPLATES_PATH $SCHEMA_PATH

# Copy template files
cp -r templates/ $CONFIG_DIR/

# Copy Schema files
cp -r schema/ $CONFIG_DIR/

# Change ownership only in case of VM/BM environment, since containers will already have access to schema and template files
if [ ! -f "/.container-env" ]; then
    chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $SCHEMA_PATH/
    chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $TEMPLATES_PATH/
fi

echo "Completed $SERVICE_NAME config upgrade to v4.0.0"
