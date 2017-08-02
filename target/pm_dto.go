package target



import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (pm *PhysicalMachine) BuildDTO() (*proto.EntityDTO, error) {
	bought, _ := pm.createCommoditiesBought(pm.ClusterID)
	sold, _ := pm.createCommoditiesSold()

	entity, err := builder.
	NewEntityDTOBuilder(proto.EntityDTO_PHYSICAL_MACHINE, pm.UUID).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		DisplayName(pm.Name).
		BuysCommodities(bought).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for pod(%v): %v",
			pm.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (pm *PhysicalMachine) createCommoditiesBought(clusterId string) (*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	clusterComm, _ := CreateKeyCommodity(clusterId, proto.CommodityDTO_CLUSTER)
	result = append(result, clusterComm)
	return result, nil
}

func (pm *PhysicalMachine) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpu := &(pm.CPU)
	cpuComm, _ := CreateResourceCommodity(cpu, proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	mem := &(pm.Memory)
	memComm, _ := CreateResourceCommodity(mem, proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	clusterComm, _ := CreateKeyCommodity(pm.UUID, proto.CommodityDTO_CLUSTER)
	result = append(result, clusterComm)

	return result, nil
}
