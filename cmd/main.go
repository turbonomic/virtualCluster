package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/containerChain/pkg/action"
	"github.com/songbinliu/containerChain/pkg/discovery"
	"github.com/songbinliu/containerChain/pkg/registration"
	"github.com/songbinliu/containerChain/pkg/target"

	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
)

var (
	targetConf   string
	opsMgrConf   string
	topologyConf string
	stitchType   string = "IP"
)

func getFlags() {
	flag.StringVar(&opsMgrConf, "turboConf", "./conf/turbo.json", "configuration file of OpsMgr")
	flag.StringVar(&targetConf, "targetConf", "./conf/target.json", "configuration file of target")
	flag.StringVar(&topologyConf, "topologyConf", "./conf/topology.conf", "topology definition of the target")

	//flag.Set("alsologtostderr", "true")
	flag.Parse()
}

func buildCluster(clusterId, clusterName, topoConf string) *target.Cluster {
	builder := target.NewClusterBuilder(clusterId, clusterName, topoConf)
	if builder == nil {
		err := fmt.Errorf("failed to create a cluster builder[%s]", topoConf)
		glog.Error(err.Error())
		return nil
	}

	cluster, err := builder.GenerateCluster()
	if err != nil {
		err := fmt.Errorf("failed to create a cluster: %v", err)
		glog.Error(err.Error())
		return nil
	}

	cluster.GenerateDTOs()
	return cluster
}

func buildProbe(stype, targetConf, topoConf string, stop chan struct{}) (*probe.ProbeBuilder, error) {

	//1. generate the target Cluster
	clusterId := "clusterId"
	clusterName := "clusterName"
	cluster := buildCluster(clusterId, clusterName, topoConf)
	if cluster == nil {
		err := fmt.Errorf("failed to build cluster[%s]", topoConf)
		glog.Error(err.Error())
		return nil, err
	}

	//2. generate clients and handlers
	config, err := discovery.NewTargetConf(targetConf)
	if err != nil {
		return nil, fmt.Errorf("failed to load json conf:%v", err.Error())
	}

	regClient := registration.NewRegistrationClient(stype)
	discoveryClient := discovery.NewDiscoveryClient(config, cluster)
	actionHandler := action.NewActionHandler(cluster, stop)

	builder := probe.NewProbeBuilder(config.TargetType, config.ProbeCategory).
		RegisteredBy(regClient).
		DiscoversTarget(config.Address, discoveryClient).
		ExecutesActionsBy(actionHandler)

	return builder, nil
}

func createTapService() (*service.TAPService, error) {
	turboConfig, err := service.ParseTurboCommunicationConfig(opsMgrConf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpsMgrConfig: %v", err)
	}

	stop := make(chan struct{})
	probeBuilder, err := buildProbe(stitchType, targetConf, topologyConf, stop)
	if err != nil {
		return nil, fmt.Errorf("failed to create probe: %v", err)
	}

	tapService, err := service.NewTAPServiceBuilder().
		WithTurboCommunicator(turboConfig).
		WithTurboProbe(probeBuilder).
		Create()

	if err != nil {
		return nil, fmt.Errorf("error when creating TapService: %v", err.Error())
	}

	return tapService, nil
}

func main() {
	getFlags()
	glog.V(2).Infof("hello")
	defer glog.V(2).Infof("bye")

	tap, err := createTapService()
	if err != nil {
		glog.Errorf("failed to create tapServier: %v", err)
	}

	tap.ConnectToTurbo()
	select {}
	//stop := make(chan struct{})
	//buildProbe("IP", targetConf, topologyConf, stop)
}
