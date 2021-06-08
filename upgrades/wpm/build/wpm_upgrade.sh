#!/bin/bash

NEW_EXEC_NAME="wpm"
CURRENT_VERSION=v4.0.0
BACKUP_PATH=${BACKUP_PATH:-"/tmp/"}
LOG_FILE=${LOG_FILE:-"/tmp/$NEW_EXEC_NAME-upgrade.log"}
INSTALLED_EXEC_PATH="/opt/$NEW_EXEC_NAME/bin/$NEW_EXEC_NAME"
CONFIG_PATH="/etc/$NEW_EXEC_NAME"
echo "" >$LOG_FILE

if [ -d "/opt/workload-policy-manager/secure-docker-daemon" ]; then
  ./upgrade-secure-docker-daemon.sh |& tee -a $LOG_FILE
fi

./upgrade.sh -v $CURRENT_VERSION -e $INSTALLED_EXEC_PATH -c $CONFIG_PATH -n $NEW_EXEC_NAME -b $BACKUP_PATH |& tee -a $LOG_FILE
exit ${PIPESTATUS[0]}