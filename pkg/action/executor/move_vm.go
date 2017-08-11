package executor

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/containerChain/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type VirtualMachineMover struct {
	cluster *target.ClusterHandler
}

func NewVirtualMachineMover(c *target.ClusterHandler) *VirtualMachineMover {
	return &VirtualMachineMover{
		cluster: c,
	}
}

func (m *VirtualMachineMover) Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) error {
	glog.V(2).Infof("begin to move a VirtualMachine.")

	//1. check
	vmEntity := actionItem.GetTargetSE()
	if vmEntity == nil {
		return fmt.Errorf("TargetSE is empty.")
	}

	hostEntity := actionItem.GetNewSE()
	if hostEntity == nil {
		return fmt.Errorf("HostEntity is empty.")
	}

	hostType := hostEntity.GetEntityType()
	if hostType != proto.EntityDTO_PHYSICAL_MACHINE {
		return fmt.Errorf("new host entity is not a VM: %v", hostType)
	}

	//2. move
	vmId := vmEntity.GetId()
	hostId := hostEntity.GetId()

	glog.V(2).Infof("move vnodeId: %s, new NodeId:%s", vmId, hostId)
	err := m.cluster.MoveVirtualMachine(vmId, hostId)
	if err != nil {
		return fmt.Errorf("move failed: %v", err)
	}

	return nil
}
