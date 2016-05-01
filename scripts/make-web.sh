#!/bin/bash
set -e
set -x

pushd web
  pushd elm
    make
  popd

  cp elm/bundle.js public/static/app/
popd
