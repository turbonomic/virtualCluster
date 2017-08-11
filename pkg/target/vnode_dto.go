package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

//  ---------- Virtual Machine Node ------------------
func (vnode *VNode) BuildDTO(pm *Node) (*proto.EntityDTO, error) {
	sold, _ := vnode.createCommoditiesSold()
	bought, _ := vnode.createCommoditiesBought()
	provider := builder.CreateProvider(proto.EntityDTO_PHYSICAL_MACHINE, pm.UUID)

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_MACHINE, vnode.UUID).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		DisplayName(vnode.Name).
		VirtualMachineData(vnode.getVMRData()).
		Provider(provider).
		BuysCommodities(bought).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for pod(%v): %v",
			vnode.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (vnode *VNode) getVMRData() *proto.EntityDTO_VirtualMachineData {
	ips := []string{vnode.IP}
	connected := true

	vmState := &proto.EntityDTO_VMState{
		Connected: &connected,
	}

	vmData := &proto.EntityDTO_VirtualMachineData{
		IpAddress: ips,
		VmState:   vmState,
		GuestName: &(vnode.Name),
	}
	return vmData
}

func (vnode *VNode) createCommoditiesBought() ([]*proto.CommodityDTO, error) {
	cpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU).Used(vnode.CPU.Capacity).Create()
	memComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM).Used(vnode.Memory.Capacity).Create()
	clusterComm, _ := CreateKeyCommodityBought(vnode.ClusterId, proto.CommodityDTO_CLUSTER)

	return []*proto.CommodityDTO{cpuComm, memComm, clusterComm}, nil
}

func (vnode *VNode) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpu := &(vnode.CPU)
	cpuComm, _ := CreateResourceCommodity(cpu, proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	mem := &(vnode.Memory)
	memComm, _ := CreateResourceCommodity(mem, proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	clusterComm, _ := CreateKeyCommodity(vnode.ClusterId, proto.CommodityDTO_CLUSTER)
	result = append(result, clusterComm)

	return result, nil
}

func (vnode *VNode) BuildSubDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	for _, pod := range vnode.Pods {
		podDTO, err := pod.BuildDTO(vnode)
		if err != nil {
			e := fmt.Errorf("failed to build PodDTO for node[%s] pod[%s]", vnode.Name, pod.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, podDTO)

		subDTOs, err := pod.BuildContainerDTOs()
		if err != nil {
			e := fmt.Errorf("failed to build Pod-containerDTOs for node[%s] pod[%s]",
				vnode.Name, pod.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, subDTOs...)
	}

	glog.V(3).Infof("There are %d DTOs for VNode[%s].", len(result)+1, vnode.Name)

	return result, nil
}
