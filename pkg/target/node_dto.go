package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

//  ---------- Virtual Machine Node ------------------
func (node *Node) BuildDTO() (*proto.EntityDTO, error) {
	//bought, _ := node.createCommoditiesBought(node.ClusterID)
	sold, _ := node.createCommoditiesSold()

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_PHYSICAL_MACHINE, node.UUID).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		DisplayName(node.Name).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for pod(%v): %v",
			node.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	node.addPMRelatedData(entity)

	return entity, nil
}

func (node *Node) addPMRelatedData(e *proto.EntityDTO) error {
	//what is Id for?
	memId := fmt.Sprintf("mem-%s", node.UUID)
	mem := &proto.EntityDTO_MemoryData{
		Id:       &memId,
		Capacity: &(node.Memory.Capacity),
	}

	cpuId := fmt.Sprintf("cpu-%s", node.UUID)
	cpu := &proto.EntityDTO_ProcessorData{
		Id:       &cpuId,
		Capacity: &(node.CPU.Capacity),
	}

	relatedData := &proto.EntityDTO_PhysicalMachineRelatedData{
		Memory:    mem,
		Processor: []*proto.EntityDTO_ProcessorData{cpu},
	}

	e.RelatedEntityData = &proto.EntityDTO_PhysicalMachineRelatedData_{PhysicalMachineRelatedData: relatedData}
	return nil
}

func (node *Node) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpu := &(node.CPU)
	cpuComm, _ := CreateResourceCommodity(cpu, proto.CommodityDTO_CPU)
	result = append(result, cpuComm)

	mem := &(node.Memory)
	memComm, _ := CreateResourceCommodity(mem, proto.CommodityDTO_MEM)
	result = append(result, memComm)

	clusterComm, _ := CreateKeyCommodity(node.ClusterId, proto.CommodityDTO_CLUSTER)
	result = append(result, clusterComm)

	return result, nil
}

// build DTOs for the hosted VNodes (Virtual Machine)
func (node *Node) BuildSubDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	for _, vm := range node.VMs {
		vmDTO, err := vm.BuildDTO(node)
		if err != nil {
			e := fmt.Errorf("failed to build VMDTO for node[%s] vnode[%s]", node.Name, vm.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, vmDTO)

		subDTOs, err := vm.BuildSubDTOs()
		if err != nil {
			e := fmt.Errorf("failed to build VM-PodDTOs for node[%s] vnode[%s]",
				node.Name, vm.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, subDTOs...)
	}

	return result, nil
}
