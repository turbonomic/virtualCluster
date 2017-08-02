package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (pod *Pod) BuildDTO(host *PhysicalMachine) (*proto.EntityDTO, error) {
	bought, _ := pod.createCommoditiesBought(host.ClusterID)
	sold, _ := pod.createCommoditiesSold()
	provider := builder.CreateProvider(proto.EntityDTO_PHYSICAL_MACHINE, host.UUID)

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_CONTAINER_POD, pod.UUID).
		DisplayName(pod.Name).
		Provider(provider).
		BuysCommodities(bought).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for pod(%v): %v",
			pod.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (pod *Pod) createCommoditiesBought(clusterId string) (*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpuComm, _ := CreateResourceCommodity(&(pod.CPU), proto.CommodityDTO_CPU)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodity(&(pod.Memory), proto.CommodityDTO_MEM)
	result = append(result, memComm)

	clusterComm, _ := CreateKeyCommodity(clusterId, proto.CommodityDTO_CLUSTER)
	result = append(result, clusterComm)
	return result, nil
}

func (pod *Pod) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO
	cpuComm, _ := CreateResourceCommodity(&(pod.CPU), proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodity(&(pod.Memory), proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	podComm, _ := CreateKeyCommodity(pod.UUID, proto.CommodityDTO_VMPM_ACCESS)
	result = append(result, podComm)

	return result, nil
}
