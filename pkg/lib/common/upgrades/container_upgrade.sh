#!/bin/bash

exit_on_error() {
  if [ $? != 0 ]; then
    echo "$1"
    exit 1
  fi
}

help() {
  echo "
  This is a generic upgrade script intended to help upgrade config and database changes of component to the latest version.
  ./hvs-container-upgrade-v*.bin {CURRENT_VERSION} {PATH_TO_CONFIG_DIR}
"
  exit 0
}

main() {
  COMPONENT_VERSION=$1
  CONFIG_DIR=$2

  if [ "$1" == "help" ] || [ "$1" == "--help" ] || [ "$1" == "-h" ] || [ "$1" == "" ] ; then
    help "$@"
  fi

  echo "Migrating Database if required"
  ./config_upgrade.sh $COMPONENT_VERSION $CONFIG_DIR "./database" ""
  exit_on_error "Failed to upgrade the database to the latest."

  echo "Migrating Configuration if required"
  ./config_upgrade.sh $COMPONENT_VERSION $CONFIG_DIR "./config" ".sh"
  exit_on_error "Failed to upgrade the configuration to the latest."
}

main "$@"
