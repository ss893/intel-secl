#!/bin/bash

SECURE_DOCKER_DAEMON_BACKUP_PATH="/tmp/workload-policy-manager_backup/container-runtime"
echo "Upgrading secure-docker-daemon"

systemctl stop docker

echo "Creating backup $SECURE_DOCKER_DAEMON_BACKUP_PATH directory for secure-docker-daemon binary"
mkdir -p $SECURE_DOCKER_DAEMON_BACKUP_PATH
echo "Taking backup of secure-docker-daemon"
cp /usr/bin/docker $SECURE_DOCKER_DAEMON_BACKUP_PATH
which /usr/bin/dockerd-ce 2>/dev/null
if [ $? -ne 0 ]; then
  cp /usr/bin/dockerd $SECURE_DOCKER_DAEMON_BACKUP_PATH
else
  cp /usr/bin/dockerd-ce $SECURE_DOCKER_DAEMON_BACKUP_PATH
fi


cp -f docker-daemon/docker /usr/bin/
which /usr/bin/dockerd-ce 2>/dev/null
if [ $? -ne 0 ]; then
  cp -f docker-daemon/dockerd-ce /usr/bin/dockerd
else
  cp -f docker-daemon/dockerd-ce /usr/bin/dockerd-ce
fi

echo "Starting secure-docker-daemon"
systemctl start docker
if [ $? -ne 0 ]; then
  echo "Error while starting secure-docker-daemon"
  exit 1
else
  echo "Upgraded secure-docker-daemon successfully"
fi