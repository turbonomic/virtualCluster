package executor

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/virtualCluster/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type PodMover struct {
	cluster *target.ClusterHandler
}

func NewPodMover(c *target.ClusterHandler) *PodMover {
	return &PodMover{
		cluster: c,
	}
}

func (m *PodMover) Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) error {
	glog.V(2).Infof("begin to move a Pod.")

	//1. check
	podEntity := actionItem.GetTargetSE()
	if podEntity == nil {
		return fmt.Errorf("TargetSE is empty.")
	}

	hostEntity := actionItem.GetNewSE()
	if hostEntity == nil {
		return fmt.Errorf("HostEntity is empty.")
	}
	hostType := hostEntity.GetEntityType()
	//if hostType != proto.EntityDTO_VIRTUAL_MACHINE && hostType != proto.EntityDTO_PHYSICAL_MACHINE {
	if hostType != proto.EntityDTO_VIRTUAL_MACHINE {
		return fmt.Errorf("new host entity is not a VM: %v", hostType)
	}

	//2. move
	podId := podEntity.GetId()
	hostId := hostEntity.GetId()

	glog.V(2).Infof("podId: %s, new VNodeId:%s", podId, hostId)
	err := m.cluster.MovePod(podId, hostId)
	if err != nil {
		return fmt.Errorf("move failed: %v", err)
	}

	return nil
}
