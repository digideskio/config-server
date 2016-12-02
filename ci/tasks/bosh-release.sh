#!/bin/sh
set -e -x

install bosh-cli/bosh-cli* /usr/local/bin/bosh

cd config-server
bosh create-release --force --tarball --name config-server --version acceptance