#!/bin/sh

## Usage:
##   . ./env-unset.sh

unset $(grep -v '^#' .env | sed -E 's/CFGSERVICE_(.*)=.*/\1/' | xargs)