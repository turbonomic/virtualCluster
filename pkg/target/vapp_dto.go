package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (vapp *VirtualApp) BuildDTO() (*proto.EntityDTO, error) {
	vAppBuilder := builder.
		NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, vapp.UUID).
		DisplayName(vapp.Name)

	if err := vapp.getCommoditiesBought(vAppBuilder); err != nil {
		nerr := fmt.Errorf("build VApplication DTO failed: %v", err)
		glog.Error(nerr.Error())
		return nil, nerr
	}

	vappData := &proto.EntityDTO_VirtualApplicationData{
		ServiceType: &(vapp.Name),
	}
	vAppBuilder.VirtualApplicationData(vappData)
	vAppBuilder.WithPowerState(proto.EntityDTO_POWERED_ON)

	entity, err := vAppBuilder.Create()
	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for vApplication(%v): %v",
			vapp.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (vapp *VirtualApp) getCommoditiesBought(vAppBuilder *builder.EntityDTOBuilder) error {
	i := 0

	for _, pod := range vapp.Pods {
		for _, container := range pod.Containers {
			if container.App == nil {
				glog.Errorf("contain.App is not ready; VirtualApp[%s]-pod[%s]-container[%s]",
					vapp.Name, pod.Name, container.Name)
				continue
			}

			app := container.App
			appCommodity, err := builder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).
				Key(app.UUID).
				Used(app.QPS.Used).
				Create()
			if err != nil {
				glog.Errorf("failed to create commodity bought for VirtualApp[%s]-pod[%s]-container[%s]-app[%s]",
					vapp.Name, pod.Name, container.Name, app.Name)
				continue
			}

			bought := []*proto.CommodityDTO{appCommodity}

			appProvider := builder.CreateProvider(proto.EntityDTO_APPLICATION, app.UUID)
			vAppBuilder.Provider(appProvider).BuysCommodities(bought)
			i++
		}
	}

	if i < 1 {
		return fmt.Errorf("cannot get commodities bought from containers for VApp[%s/%s]", vapp.Kind, vapp.Name)
	}

	return nil
}

/*
func (vapp *VirtualApp) getCommodityBought(appDTO *proto.EntityDTO)([]*proto.CommodityDTO, error) {
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
}*/
