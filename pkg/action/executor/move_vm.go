package executor

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/containerChain/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type VirtualMachineMover struct {
	cluster *target.Cluster
}

func NewVirtualMachineMover(c *target.Cluster) *VirtualMachineMover {
	return &VirtualMachineMover{
		cluster: c,
	}
}

func (m *VirtualMachineMover) Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) error {
	glog.V(2).Infof("begin to move a VirtualMachine.")
	return fmt.Errorf("VirtualMachiner Mover is not implemented.")
}
