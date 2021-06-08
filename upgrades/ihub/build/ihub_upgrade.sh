#!/bin/bash

CURRENT_VERSION=v4.0.0
BACKUP_PATH=${BACKUP_PATH:-"/tmp/"}
BACKUP_PATH=${BACKUP_PATH}ihub/
INSTALLED_EXEC_PATH=/opt/ihub/bin/ihub
NEW_EXEC_NAME=ihub

IHUB_INSTANCES=( $(pgrep -a ihub | grep "/usr/bin/ihub" | awk 'NF>1{print $NF}') )
declare -p IHUB_INSTANCES &>/dev/null

INDEX=0
for i in "${IHUB_INSTANCES[@]}"
do
  echo "Upgrading iHUB instance - $i"
  INSTANCE_NAME=${i}
  LOG_FILE=${LOG_FILE:-"/tmp/${INSTANCE_NAME}-upgrade.log"}
  CONFIG_PATH="/etc/${INSTANCE_NAME}"
  echo "" > $LOG_FILE
  if [ ${INDEX} -eq 1 ] ; then
    export BACKUP_ONLY=true
  fi
  INDEX=$(expr $INDEX + 1)
  ihub stop -i $INSTANCE_NAME
  ./upgrade.sh -v $CURRENT_VERSION -b $BACKUP_PATH -e $INSTALLED_EXEC_PATH -n $NEW_EXEC_NAME -c $CONFIG_PATH |& tee -a $LOG_FILE
  if [ ${PIPESTATUS[0]} != 0 ]; then
    exit ${PIPESTATUS[0]}
  fi
  ihub start -i $INSTANCE_NAME
  echo "Successfully upgraded ${INSTANCE_NAME} instance of iHUB"
done