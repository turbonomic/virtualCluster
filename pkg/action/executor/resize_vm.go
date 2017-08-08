package executor

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/songbinliu/containerChain/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type VMResizer struct {
	cluster *target.ClusterHandler
}

func NewVMResizer(c *target.ClusterHandler) *VMResizer {
	return &VMResizer{
		cluster: c,
	}
}

func (m *VMResizer) Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) error {
	glog.V(2).Infof("begin to move a Pod.")
	return fmt.Errorf("VirtualMachine Resizer is not implemented.")
}
