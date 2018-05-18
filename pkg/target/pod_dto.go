package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (pod *Pod) BuildDTO(host *VNode) (*proto.EntityDTO, error) {
	bought, _ := pod.createCommoditiesBought(host.ClusterId)
	sold, _ := pod.createCommoditiesSold()
	provider := builder.CreateProvider(proto.EntityDTO_PHYSICAL_MACHINE, host.UUID)

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_CONTAINER_POD, pod.UUID).
		DisplayName(pod.Name).
		Provider(provider).
		BuysCommodities(bought).
		SellsCommodities(sold).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for pod(%v): %v",
			pod.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (pod *Pod) createCommoditiesBought(clusterId string) ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpuComm, _ := CreateResourceCommodityBought(&(pod.CPU), proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodityBought(&(pod.Memory), proto.CommodityDTO_VMEM)
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

func (pod *Pod) createContainerPodData() *proto.EntityDTO_ContainerPodData {
	// Add IP address in ContainerPodData. Some pods (system pods and daemonset pods) may use the host IP as the pod IP,
	// in which case the IP address will not be unique (in the k8s cluster) and hence not populated in ContainerPodData.
	fullName := pod.Name
	ns := "ns"
	return &proto.EntityDTO_ContainerPodData{
		// Note the port needs to be set if needed
		IpAddress: &fullName,
		FullName:  &fullName,
		Namespace: &ns,
	}
}

func (pod *Pod) BuildContainerDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	for _, container := range pod.Containers {
		containerDTO, err := container.BuildDTO(pod)
		if err != nil {
			e := fmt.Errorf("failed to build containerDTO for pod[%s] container[%s]",
				pod.Name, container.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, containerDTO)

		appDTO, err := container.BuildAppDTO(pod)
		if err != nil {
			e := fmt.Errorf("failed to build appDTO for pod[%s] container[%s]",
				pod.Name, container.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, appDTO)
	}

	glog.V(3).Infof("There are %d DTOs for Pod[%s].", len(result)+1, pod.Name)

	return result, nil
}
