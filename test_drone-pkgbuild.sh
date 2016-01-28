#!/bin/bash

# package to build
pkg=${1:-"neovim-git"}

mkdir drone_dir

docker run -v $(pwd)/drone_dir:/drone -i mikkeloscar/drone-pkgbuild <<EOF
{
    "system": {
        "link_url": "https://drone.server.com/api/?repos/mikkeloscar/test/builds/1?fork=true&access_token=na&pkgs=$pkg&src=aur"
    },
    "workspace": {
        "path": "/drone",
        "root": "",
        "keys": {
            "private": "",
            "public": ""
        }
    },
    "vargs": {
        "sign_key": "test key"
    }
}
EOF

rm -rf drone_dir
