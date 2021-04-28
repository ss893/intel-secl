#!/bin/bash

NEW_EXEC_NAME="wpm"
CURRENT_VERSION=v3.6.0
BACKUP_PATH=${BACKUP_PATH:-"/tmp/"}
INSTALLED_EXEC_PATH="/opt/workload-policy-manager/bin/$NEW_EXEC_NAME"
CONFIG_PATH="/etc/workload-policy-manager"
OLD_EXEC_NAME="workload-policy-manager"
LOG_FILE=${LOG_FILE:-"/tmp/$OLD_EXEC_NAME-upgrade.log"}
echo "" > $LOG_FILE
./upgrade.sh -v $CURRENT_VERSION -e $INSTALLED_EXEC_PATH -c $CONFIG_PATH -n $NEW_EXEC_NAME -o $OLD_EXEC_NAME -b $BACKUP_PATH |& tee -a $LOG_FILE
