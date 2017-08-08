package executor

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/containerChain/pkg/target"

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
	glog.V(2).Infof("begin to move a Pod.")
	return fmt.Errorf("Container Resizer is not implemented.")
}
