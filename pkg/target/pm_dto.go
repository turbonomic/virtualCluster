package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (node *HostNode) BuildDTO() (*proto.EntityDTO, error) {
	//bought, _ := node.createCommoditiesBought(node.ClusterID)
	sold, _ := node.createCommoditiesSold()

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_MACHINE, node.UUID).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		DisplayName(node.Name).
		VirtualMachineData(node.getVMRData()).
		//BuysCommodities(bought).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for pod(%v): %v",
			node.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	//node.addPMRelatedData(entity)

	return entity, nil
}

func (node *HostNode) getVMRData() *proto.EntityDTO_VirtualMachineData {
	ips := []string{node.IP}
	connected := true

	vmState := &proto.EntityDTO_VMState{
		Connected: &connected,
	}

	vmData := &proto.EntityDTO_VirtualMachineData{
		IpAddress: ips,
		VmState: vmState,
		GuestName: &(node.Name),
	}
	return vmData
}

func (pm *HostNode) addPMRelatedData(e *proto.EntityDTO) error {
	mem := &proto.EntityDTO_MemoryData{
		Capacity: &(pm.Memory.Capacity),
	}

	cpu := &proto.EntityDTO_ProcessorData{
		Capacity: &(pm.CPU.Capacity),
	}

	relatedData := &proto.EntityDTO_PhysicalMachineRelatedData{
		Memory: mem,
		Processor: []*proto.EntityDTO_ProcessorData{cpu},
	}

	e.RelatedEntityData = &proto.EntityDTO_PhysicalMachineRelatedData_{relatedData}
	return nil
}

//func (pm *HostNode) createCommoditiesBought(clusterId string) ([]*proto.CommodityDTO, error) {
//
//	var result []*proto.CommodityDTO
//
//	clusterComm, _ := CreateKeyCommodity(clusterId, proto.CommodityDTO_CLUSTER)
//	result = append(result, clusterComm)
//	return result, nil
//}

func (pm *HostNode) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpu := &(pm.CPU)
	cpuComm, _ := CreateResourceCommodity(cpu, proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	mem := &(pm.Memory)
	memComm, _ := CreateResourceCommodity(mem, proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	clusterComm, _ := CreateKeyCommodity(pm.ClusterID, proto.CommodityDTO_CLUSTER)
	result = append(result, clusterComm)

	return result, nil
}

func (node *HostNode) BuildPodDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	for _, pod := range node.Pods {
		podDTO, err := pod.BuildDTO(node)
		if err != nil {
			e := fmt.Errorf("failed to build PodDTO for node[%s] pod[%s]", node.Name, pod.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, podDTO)

		subDTOs, err := pod.BuildContainerDTOs()
		if err != nil {
			e := fmt.Errorf("failed to build Pod-containerDTOs for node[%s] pod[%s]",
				node.Name, pod.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, subDTOs...)
	}

	return result, nil
}
