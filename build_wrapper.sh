#!/bin/bash

config=$2

# get path from input
path=$(echo $config | jq -r '.workspace.path')

# make workspace dirs for the build and add correct permissions.
mkdir -p $path/drone_pkgbuild/{repo,sources}
chown $UGID:$UGID -R $path/drone_pkgbuild

# Run real program as user $UGNAME
echo $config | sudo -u $UGNAME /usr/bin/drone-pkgbuild
