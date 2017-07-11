package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/pflag"

	"containerChain/action"
	"containerChain/discovery"
	"containerChain/registration"

	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
)

var (
	targetConf string
	opsMgrConf string
	pmNum      int
	vmNum      int
	podNum     int
	stitchType string
)

func getFlags() {
	pflag.StringVar(&targetConf, "targetConf", "./target-conf.json", "configuration file of target")
	pflag.StringVar(&opsMgrConf, "opsMgrConf", "./turbo-conf.json", "configuration file of OpsMgr")
	pflag.IntVar(&pmNum, "pmNum", 10, "number of total physical machines")
	pflag.IntVar(&vmNum, "vmNum", 50, "number of total virtual machines")
	pflag.IntVar(&podNum, "podNum", 100, "number of total pods")
	pflag.StringVar(&stitchType, "stitchType", "IP", "stitching type (IP | UUID)")

	pflag.Set("alsologtostderr", "true")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
}

type TargetTopoConf struct {
	pmNum  int
	vmNum  int
	podNum int
}

func buildProber(stype, targetConf string, topo *TargetTopoConf, stop chan struct{}) (*probe.ProbeBuilder, error) {
	config, err := discovery.NewTargetConf(targetConf)
	if err != nil {
		return nil, fmt.Errorf("failed to load json conf:%v", err.Error())
	}

	regClient := registration.NewRegistrationClient(stype)
	discoveryClient := discovery.NewDiscoveryClient(config)
	actionHandler := action.NewActionHandler(stop)

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
	topo := TargetTopoConf{
		pmNum:  pmNum,
		vmNum:  vmNum,
		podNum: podNum,
	}

	probeBuilder, err := buildProber(stitchType, targetConf, &topo, stop)
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
	fmt.Println("vim-go")
	glog.V(2).Infof("hello")
	defer glog.V(2).Infof("bye")
}
