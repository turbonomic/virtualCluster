package target

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	defaultInfiniteCapacity = 1E10
)

func CreateCommodityBoughtWithReservation(res *Resource, reserved float64, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.NewCommodityDTOBuilder(ctype).
		Used(res.Used).
	    Reservation(reserved).
	    Resizable(true).
		Create()
}

func CreateResourceCommodityBought(res *Resource, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.NewCommodityDTOBuilder(ctype).
		Used(res.Used).
		Create()
}

func CreateResourceCommodity(res *Resource, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.NewCommodityDTOBuilder(ctype).
		Capacity(res.Capacity).
		Used(res.Used).
		Create()
}

func CreateResourceCommodityResize(res *Resource, ctype proto.CommodityDTO_CommodityType, resizeable bool) (*proto.CommodityDTO, error) {
	return builder.NewCommodityDTOBuilder(ctype).
		Capacity(res.Capacity).
		Used(res.Used).
		Resizable(resizeable).
		Create()
}

func CreateKeyCommodity(key string, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.
		NewCommodityDTOBuilder(ctype).
		Key(key).
		Capacity(defaultInfiniteCapacity).
		Create()
}

func CreateKeyCommodityBought(key string, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.
		NewCommodityDTOBuilder(ctype).
		Key(key).
		Create()
}

func CreateTransactionCommodity(key string, qps *Resource, ctype proto.CommodityDTO_CommodityType) (*proto.CommodityDTO, error) {
	return builder.
		NewCommodityDTOBuilder(ctype).
		Key(key).
		Capacity(qps.Capacity).
		Used(qps.Used).
		Create()
}
