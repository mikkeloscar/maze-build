#!/bin/bash

# Simple wrapper around maze-build to setup build directories and drop
# permissions to the build user.

# set path to pwd
path=$(pwd)

# make workspace dirs for the build and add correct permissions.
mkdir -p $path/drone_pkgbuild/{repo,sources}
chown $UGID:$UGID -R $path/drone_pkgbuild

# Run real program as user $UGNAME
sudo PLUGIN_REPO=$PLUGIN_REPO DRONE_WORKSPACE=$DRONE_WORKSPACE -u $UGNAME /usr/bin/maze-build
