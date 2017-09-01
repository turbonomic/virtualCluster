package executor

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/virtualCluster/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type ContainerResizer struct {
	cluster *target.ClusterHandler
}

func NewContainerResizer(c *target.ClusterHandler) *ContainerResizer {
	return &ContainerResizer{
		cluster: c,
	}
}

func (m *ContainerResizer) Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) error {
	glog.V(2).Infof("begin to resize a container.")
	containerSE := actionItem.GetTargetSE()
	podSE := actionItem.GetHostedBySE()
	comm := actionItem.GetNewComm()
	glog.V(2).Infof("begin to resize container[%s] hosted by pod[%s]\n comm:%++v",
		              containerSE.GetDisplayName(),
		              podSE.GetDisplayName(),
					  comm)

	cpu := -1.0
	mem := -1.0

	ctype := comm.GetCommodityType()
	switch ctype {
	case proto.CommodityDTO_VMEM:
		mem = comm.GetCapacity()
	case proto.CommodityDTO_VCPU:
		cpu = comm.GetCapacity()
	default:
		glog.Errorf("unable to resize commodity type[%v] for container[%s].", ctype, containerSE.GetId())
		return fmt.Errorf("unsupported commdity type [%v]", ctype)
	}

	if cpu < 0 && mem < 0 {
		err := fmt.Errorf("wrong new capacity: mem=%.1f, cpu=%.1f", cpu, mem)
		glog.Error(err)
		return fmt.Errorf("wrong new capacity.")
	}

	err := m.cluster.ResizeContainerCapacity(containerSE.GetId(), cpu, mem)
	if err != nil {
		glog.Errorf("Failed to resize container[%s] capacity: %v", containerSE.GetId(), err)
		return fmt.Errorf("failed to resize container capacity.")
	}

	return nil
}
