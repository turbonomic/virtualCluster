package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (app *Application) BuildDTO(docker *Container) (*proto.EntityDTO, error) {
	bought, _ := app.createCommoditiesBought(docker.UUID)
	sold, _ := app.createCommoditiesSold()
	provider := builder.CreateProvider(proto.EntityDTO_CONTAINER_POD, docker.UUID)

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_APPLICATION, app.UUID).
		DisplayName(app.Name).
		Provider(provider).
		BuysCommodities(bought).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for Application(%v/%v): %v",
			app.Name, docker.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (app *Application) createCommoditiesBought(containerId string) ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpuComm, _ := CreateResourceCommodityBought(&(app.CPU), proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodityBought(&(app.Memory), proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	podComm, _ := CreateKeyCommodity(containerId, proto.CommodityDTO_APPLICATION)
	result = append(result, podComm)
	return result, nil
}

func (app *Application) createCommoditiesSold() ([]*proto.CommodityDTO, error) {
	var result []*proto.CommodityDTO
	appComm, _ := CreateTransactionCommodity(app.UUID, proto.CommodityDTO_TRANSACTION)
	result = append(result, appComm)

	return result, nil
}
