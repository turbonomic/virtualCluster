package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (vapp *VApplication) BuildDTO(pods []*Pod,  appDTOs map[string]*proto.EntityDTO) (*proto.EntityDTO, error) {
	vAppBuilder := builder.
		NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, vapp.UUID).
		DisplayName(vapp.Name)

	if err := vapp.getCommoditiesBought(vAppBuilder, pods,appDTOs); err != nil {
		nerr := fmt.Errorf("build VApplication DTO failed: %v", err)
		glog.Error(nerr.Error())
		return nil, nerr
	}

	entity, err := vAppBuilder.Create()
	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for vApplication(%v): %v",
			vapp.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (vapp *VApplication) getCommoditiesBought(vAppBuilder *builder.EntityDTOBuilder, pods []*Pod, appDTOs map[string]*proto.EntityDTO) error {
	i := 0

	for _, pod := range pods {
		for _, container := range pod.Containers {
			appId := container.UUID

			appDTO, exist := appDTOs[appId]
			if !exist {
				//TODO: should I retrun error?
				glog.Warningf("cannot find container[%s/%s] appDTO for VApplication[%s/%s]", container.Name, appId, vapp.Kind, vapp.Name)
				continue
			}

			bought, err := vapp.getCommodityBought(appDTO)
			if err != nil {
				glog.Errorf("failed to get commodity from container[%s/%s] for VApplication[%s/%s]", container.Name, appId, vapp.Kind, vapp.Name)
				continue
			}

			appProvider := builder.CreateProvider(proto.EntityDTO_APPLICATION, appId)
			vAppBuilder.Provider(appProvider).BuysCommodities(bought)
			i ++
		}
	}

	if i < 1 {
		return fmt.Errorf("cannot get commodities bought from containers for VApp[%s/%s]", vapp.Kind, vapp.Name)
	}

	return nil
}

func (vapp *VApplication) getCommodityBought(appDTO *proto.EntityDTO)([]*proto.CommodityDTO, error) {
	var result []*proto.CommodityDTO

	hmap := make(map[proto.CommodityDTO_CommodityType]struct{})
	hmap[proto.EntityDTO_APPLICATION] = struct{}{}

	sold := appDTO.GetCommoditiesSold()
	for _, comm := range sold {
		if _, exist := hmap[comm.GetCommodityType()]; !exist {
			continue
		}

		commBought, err := builder.NewCommodityDTOBuilder(comm.GetCommodityType()).
		                      Key(comm.GetKey()).
							  Used(comm.GetUsed()).
		                      Create()
		if err != nil {
			glog.Warning("failed to create commodity.")
			continue
		}

		result = append(result, commBought)
	}

	if len(result) < 1 {
		err := fmt.Errorf("cannot find commdity for VApplication.")
		glog.Error(err.Error())
		return result, err
	}

	return result, nil
}
