package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/virtualCluster/pkg/action"
	"github.com/songbinliu/virtualCluster/pkg/discovery"
	"github.com/songbinliu/virtualCluster/pkg/registration"
	"github.com/songbinliu/virtualCluster/pkg/target"
	"github.com/songbinliu/virtualCluster/pkg/topology"

	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
)

var (
	targetConf   string
	opsMgrConf   string
	topologyConf string
	stitchType   string = "IP"
	clusterName  string = "clusterName-1"
	clusterId    string = "clusterId-1"
)

func getFlags() {
	flag.StringVar(&opsMgrConf, "turboConf", "./conf/turbo.json", "configuration file of OpsMgr")
	flag.StringVar(&targetConf, "targetConf", "./conf/target.json", "configuration file of target")
	flag.StringVar(&topologyConf, "topologyConf", "./conf/topology.conf", "topology definition of the target")
	flag.StringVar(&clusterName, "clusterName", "clusterName-1", "virtual cluster Name")
	flag.StringVar(&clusterId, "clusterId", "clusterId-1", "virtual cluster Id")

	//flag.Set("alsologtostderr", "true")
	flag.Parse()
}

func buildCluster(clusterId, clusterName, topoConf string) *target.Cluster {
	builder := topology.NewClusterBuilder(clusterId, clusterName, topoConf)
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

	//TODO: delete it, this is for test only.
	cluster.GenerateDTOs()
	return cluster
}

func buildClusterHandler(topoConf string) (*target.ClusterHandler, error) {
	cluster := buildCluster(clusterId, clusterName, topoConf)
	if cluster == nil {
		err := fmt.Errorf("failed to build cluster[%s]", topoConf)
		glog.Error(err.Error())
		return nil, err
	}

	handler := target.NewClusterHandler(cluster)
	return handler, nil
}

func buildProbe(stype, targetConf, topoConf string, stop chan struct{}) (*probe.ProbeBuilder, error) {

	//1. generate the target Cluster Handler
	clusterHandler, err := buildClusterHandler(topologyConf)
	if err != nil {
		err := fmt.Errorf("failed to build cluster handler for [%s]", topoConf)
		glog.Error(err.Error())
		return nil, err
	}

	//2. generate clients and handlers
	config, err := discovery.NewTargetConf(targetConf)
	if err != nil {
		return nil, fmt.Errorf("failed to load json conf:%v", err.Error())
	}

	regClient := registration.NewRegClient(stype)
	discoveryClient := discovery.NewDiscoveryClient(config, clusterHandler)
	actionHandler := action.NewActionHandler(clusterHandler, stop)

	builder := probe.NewProbeBuilder(config.TargetType, config.ProbeCategory).
		RegisteredBy(regClient).
		WithActionPolicies(regClient).
		WithEntityMetadata(regClient).
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

	tap, err := createTapService()
	if err != nil {
		glog.Errorf("failed to create tapServier: %v", err)
	}

	tap.ConnectToTurbo()
}
