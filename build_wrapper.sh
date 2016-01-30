#!/bin/bash

# Simple wrapper around maze-build to setup build directories and drop
# permissions to the build user.

build_vars=$2

# get path from input
path=$(echo $build_vars | jq -r '.workspace.path')

# make workspace dirs for the build and add correct permissions.
mkdir -p $path/drone_pkgbuild/{repo,sources}
chown $UGID:$UGID -R $path/drone_pkgbuild

# Run real program as user $UGNAME
echo $build_vars | sudo -u $UGNAME /usr/bin/maze-build
