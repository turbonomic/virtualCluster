#!/bin/bash



options="--logtostderr"
options="$options --v 3"
options="$options --topologyConf ./conf/topology.conf"
options="$options --turboConf ./conf/turbo.json"
options="$options --targetConf ./conf/target.json"
options="$options --clusterName myCluster"

set -x

#1. build it
make build
ret=$?
if [ $ret -ne 0 ] ; then
    echo "Build failed."
    exit 1
fi

#2. run it
./_output/vCluster $options
