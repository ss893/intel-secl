#!/bin/bash
#iterate over all files in "config" directory which installer must copy
# run all config upgrade scripts from current version to the latest version
# e.g. v3.5.0 to v3.7.0 upgrade should run all following scripts
# v3.5.1_{comopnent}_config.sh
# v3.6.0_{comopnent}_config.sh
# v3.6.1_{comopnent}_config.sh
# v3.6.2_{comopnent}_config.sh
# v3.7.0_{comopnent}_config.sh
# but config update script should not run anything like
# v3.2.0_{comopnent}_config.sh
# v3.4.0_{comopnent}_config.sh
# v3.4.1_{comopnent}_config.sh
# v3.5.0_{comopnent}_config.sh

#Upgrade config
echo "Config upgrade started"
#get currently installed version number after removing '.'
COMPONENT_VERSION=$(echo $1 | sed 's/v//' | sed 's/\.//g')
READ_FILES=false
CONFIG_DIR="./config"
if [ -d "$CONFIG_DIR" ]; then
  chmod +x $CONFIG_DIR/*.sh
  #Sort files
  cd $CONFIG_DIR && ls -1 *.sh | sort -n -k1.4 >temp_configs
  IFS=$'\r\n' GLOBIGNORE='*' command eval 'configUpgradeFiles=($(cat temp_configs))'
  rm -rf temp_configs

  for i in "${configUpgradeFiles[@]}"; do
    :
    #get script version number after removing '.'
    VERSION=$(echo $i | cut -d'_' -f1 | sed 's/v//' | sed 's/\.//g')
    #Ignore files till component version is matched
    if [ $VERSION -gt $COMPONENT_VERSION ]; then
      READ_FILES=true
    fi

    #Run all config files which are post current release
    if $READ_FILES; then
      echo "Running upgrade script - $i with arguments $2"
      ./$i $2
      if [ $? != 0 ]; then
        echo "Failed to apply $i upgrade script"
        exit 1
      fi
    fi
  done
  cd -
fi
if $READ_FILES; then
  echo "Config upgraded successfully"
else
  echo "Existing config is good for the latest release"
fi
