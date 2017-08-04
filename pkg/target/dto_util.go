package target

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
)

const (
	defaultInfiniteCapacity = 1E10
)

func CreateResourceCommodity(res *Resource, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.NewCommodityDTOBuilder(ctype).
		Capacity(res.Capacity).
		Used(res.Used).
		Create()
}

func CreateKeyCommodity(key string, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.
		NewCommodityDTOBuilder(ctype).
		Key(key).
		Capacity(defaultInfiniteCapacity).
		Create()
}
