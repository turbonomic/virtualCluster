#!/bin/bash



options="--logtostderr"
options="$options --v 3"
options="$options --topologyConf ./conf/test.conf"
options="$options --turboConf ./conf/turbo.json"
options="$options --targetConf ./conf/target.json"

set -x

#1. build it
make build
ret=$?
if [ $ret -ne 0 ] ; then
    echo "Build failed."
    exit 1
fi

#2. run it
./_output/containerChain $options
