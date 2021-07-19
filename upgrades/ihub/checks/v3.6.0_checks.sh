#!/bin/bash

if [[ -z $HVS_BASE_URL && -z $SHVS_BASE_URL ]] ; then
  echo "HVS_BASE_URL 0r SHVS_BASE_URL is required for the upgrade to v3.6.0"
  exit 1
fi
