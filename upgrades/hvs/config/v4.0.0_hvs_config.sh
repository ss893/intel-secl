#!/bin/bash

SERVICE_NAME=hvs
SERVICE_USERNAME=hvs
CONFIG_DIR="/etc/$SERVICE_NAME"
CONFIG_FILE="$CONFIG_DIR/config.yml"

echo "Starting $SERVICE_NAME config upgrade to v4.0.0"
TEMPLATES_PATH=$CONFIG_DIR/templates

mkdir -p $TEMPLATES_PATH
cp -r templates/ $CONFIG_DIR/ && chown $SERVICE_USERNAME:$SERVICE_USERNAME $TEMPLATES_PATH/*
$SERVICE_NAME setup create-default-flavor-template

echo "Completed $SERVICE_NAME config upgrade to v4.0.0"
