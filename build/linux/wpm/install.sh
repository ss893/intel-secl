#!/bin/bash

# Check OS
OS=$(cat /etc/os-release | grep ^ID= | cut -d'=' -f2)
temp="${OS%\"}"
temp="${temp#\"}"
OS="$temp"

COMPONENT_NAME=wpm
SERVICE_ENV=wpm.env

# Upgrade if component is already installed
if command -v $COMPONENT_NAME &>/dev/null; then
  n=0
  until [ "$n" -ge 3 ]
  do
  echo "$COMPONENT_NAME is already installed, Do you want to proceed with the upgrade? [y/n]"
  read UPGRADE_NEEDED
  if [ $UPGRADE_NEEDED == "y" ] || [ $UPGRADE_NEEDED == "Y" ] ; then
    echo "Proceeding with the upgrade.."
    ./${COMPONENT_NAME}_upgrade.sh
    exit $?
  elif [ $UPGRADE_NEEDED == "n" ] || [ $UPGRADE_NEEDED == "N" ] ; then
    echo "Exiting the installation.."
    exit 0
  fi
  n=$((n+1))
  done
  echo "Exiting the installation.."
  exit 0
fi

# find .env file
echo PWD IS $(pwd)
if [ -f ~/$SERVICE_ENV ]; then
  echo Reading Installation options from $(realpath ~/$SERVICE_ENV)
  env_file=~/$SERVICE_ENV
elif [ -f ../$SERVICE_ENV ]; then
  echo Reading Installation options from $(realpath ../$SERVICE_ENV)
  env_file=../$SERVICE_ENV
fi

if [[ $EUID -ne 0 ]]; then
  echo "This installer must be run as root"
  exit 1
fi

if [ -z $env_file ]; then
  echo "No .env file found"
  WPM_NOSETUP="true"
else
  source $env_file
  env_file_exports=$(cat $env_file | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
  if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
fi

echo "Installing Workload Policy Manager..."

# default settings
PRODUCT_HOME=/opt/$COMPONENT_NAME
BIN_PATH=$PRODUCT_HOME/bin
LOG_PATH=/var/log/$COMPONENT_NAME/
CONFIG_PATH=/etc/$COMPONENT_NAME/
CERTS_PATH=$CONFIG_PATH/certs
CERTDIR_TRUSTEDCAS=$CERTS_PATH/trustedca
CERTDIR_FLAVOR_SIGN_DIR=$CERTS_PATH/flavorsign
CERTDIR_KBS_ENVELOPKEY_DIR=$CERTS_PATH/kbs
FLAVORS=$PRODUCT_HOME/flavors
VM_IMAGES_PATH=$PRODUCT_HOME/vm-images/
ENCRYPTED_VM_IMAGES_PATH=$PRODUCT_HOME/encrypted-vm-images/

for directory in $BIN_PATH $LOG_PATH $CONFIG_PATH $CERTS_PATH $CERTDIR_TRUSTEDCAS $CERTDIR_FLAVOR_SIGN_DIR $CERTDIR_KBS_ENVELOPKEY_DIR $FLAVORS $VM_IMAGES_PATH $ENCRYPTED_VM_IMAGES_PATH; do
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
  chmod 700 $directory
done

cp $COMPONENT_NAME $BIN_PATH/
chmod 700 $BIN_PATH/*
ln -sfT $BIN_PATH/$COMPONENT_NAME /usr/bin/$COMPONENT_NAME

# log file permission change
chmod 740 $LOG_PATH

auto_install() {
  local component=${1}
  local cprefix=${2}
  local packages=$(eval "echo \$${cprefix}_PACKAGES")
  # detect available package management tools. start with the less likely ones to differentiate.
  if [ "$OS" == "rhel" ]; then
    yum -y install $packages
  elif [ "$OS" == "ubuntu" ]; then
    apt -y install $packages
  fi
}

# SCRIPT EXECUTION
logRotate_clear() {
  logrotate=""
}

logRotate_detect() {
  local logrotaterc=$(ls -1 /etc/logrotate.conf 2>/dev/null | tail -n 1)
  logrotate=$(which logrotate 2>/dev/null)
  if [ -z "$logrotate" ] && [ -f "/usr/sbin/logrotate" ]; then
    logrotate="/usr/sbin/logrotate"
  fi
}

logRotate_install() {
  LOGROTATE_PACKAGES="logrotate"
  if [ "$(whoami)" == "root" ]; then
    auto_install "Log Rotate" "LOGROTATE"
    if [ $? -ne 0 ]; then
      echo "Failed to install logrotate"
      exit -1
    fi
  fi
  logRotate_clear
  logRotate_detect
  if [ -z "$logrotate" ]; then
    echo "logrotate is not installed"
  else
    echo "logrotate installed in $logrotate"
  fi
}

logRotate_install

export LOG_ROTATION_PERIOD=${LOG_ROTATION_PERIOD:-weekly}
export LOG_COMPRESS=${LOG_COMPRESS:-compress}
export LOG_DELAYCOMPRESS=${LOG_DELAYCOMPRESS:-delaycompress}
export LOG_COPYTRUNCATE=${LOG_COPYTRUNCATE:-copytruncate}
export LOG_SIZE=${LOG_SIZE:-100M}
export LOG_OLD=${LOG_OLD:-12}

mkdir -p /etc/logrotate.d

if [ ! -a /etc/logrotate.d/${COMPONENT_NAME} ]; then
  echo "/var/log/${COMPONENT_NAME}/*.log {
    missingok
    notifempty
    rotate $LOG_OLD
    maxsize $LOG_SIZE
    nodateext
    $LOG_ROTATION_PERIOD
    $LOG_COMPRESS
    $LOG_DELAYCOMPRESS
    $LOG_COPYTRUNCATE
}" >/etc/logrotate.d/${COMPONENT_NAME}
fi

# check if WPM_NOSETUP is defined
if [ "${WPM_NOSETUP,,}" == "true" ]; then
  echo "WPM_NOSETUP is true, skipping setup"
  echo "Run \"$COMPONENT_NAME setup all\" for manual setup"
  echo "Installation completed successfully!"
else
  $COMPONENT_NAME setup all --force
  SETUPRESULT=$?
  if [ ${SETUPRESULT} == 0 ]; then
    echo "Installation completed successfully!"
  else
    echo "Installation completed with errors"
  fi
fi

if [ "$WPM_WITH_CONTAINER_SECURITY_CRIO" = "y" ] || [ "$WPM_WITH_CONTAINER_SECURITY_CRIO" = "Y" ] || [ "$WPM_WITH_CONTAINER_SECURITY_CRIO" = "yes" ]; then
  isinstalled=$(rpm -q skopeo)
  if [ "$isinstalled" == "package skopeo is not installed" ]; then
    echo "Prerequisite skopeo is not installed, please install skopeo before proceeding with container confidentiality."
  fi
fi
echo "Installation completed."
