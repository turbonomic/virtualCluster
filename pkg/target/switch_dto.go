package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

//  ---------- Network Switch ------------------
func (networkswitch *Switch) BuildDTO() (*proto.EntityDTO, error) {
	//bought, _ := node.createCommoditiesBought(node.ClusterID)
	sold, _ := networkswitch.createCommoditiesSold()

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_SWITCH, networkswitch.UUID).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		DisplayName(networkswitch.Name).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for switch(%v): %v",
			networkswitch.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (networkswitch *Switch) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	net := &(networkswitch.NetworkThroughput)
	netComm, _ := CreateResourceCommodity(net, proto.CommodityDTO_NET_THROUGHPUT)
	result = append(result, netComm)

	return result, nil
}

// build DTOs for the hosted Nodes (Physical Machine)
func (networkswitch *Switch) BuildSubDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	for _, pm := range networkswitch.PMs {
		pmDTO, err := pm.BuildDTO(networkswitch)
		if err != nil {
			e := fmt.Errorf("failed to build VMDTO for networkswitch[%s] node[%s]", networkswitch.Name, pm.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, pmDTO)

		subDTOs, err := pm.BuildSubDTOs()
		if err != nil {
			e := fmt.Errorf("failed to build VM-PodDTOs for networkswitch[%s] node[%s]",
				networkswitch.Name, pm.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, subDTOs...)
	}

	return result, nil
}
