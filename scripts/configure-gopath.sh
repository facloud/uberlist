#!/bin/bash
set -e
set -x
ROOT_PATH=$(cd $(dirname $BASH_SOURCE)/..; pwd)

function is_osx {
	[ $(uname -s) == "Darwin" ]
}

function is_linux {
	[ $(uname -s) == "Linux" ]
}

function universal_ln {
  is_osx && hln $@
  is_linux && ln $@
}

[ -z $GOPATH ] && echo "No GOPATH set, bailing out.." && exit 1
is_osx && ! which hln > /dev/null 2> /dev/null && echo "Please install OSX hardlinks - https://github.com/selkhateeb/hardlink" && exit 1

srcPath=$ROOT_PATH/server
destPath=$GOPATH/src/github.com/glestaris/uberlist-server
if ! [ -d $destPath ]; then
  mkdir -p $(dirname $destPath)
  universal_ln $srcPath $destPath
fi
