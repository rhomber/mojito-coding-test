package handler

import (
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/service"
	"mojito-coding-test/common/chttp"
	"net/http"
)

func GetAuctionLotBids(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	err := ctx.GetDb().Transaction(func(txn *gorm.DB) error {
		sm := ctx.GetServiceManager().(*service.Manager)

		auctionLotId, err := ctx.URLParamInt("auctionLotId")
		if err != nil {
			return err
		}

		listDTO, err := sm.AuctionLotBid.List(txn, uint(auctionLotId))
		if err != nil {
			return err
		}

		ctx.Respond(200, listDTO)

		return nil
	})
	if err != nil {
		ctx.InternalError(err)
	}
}

func PostAuctionLotBid(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	var itemDTO dto.CreateAuctionLotBid
	if !ctx.MustDecode(&itemDTO) {
		return
	}

	err := ctx.GetDb().Transaction(func(txn *gorm.DB) error {
		sm := ctx.GetServiceManager().(*service.Manager)

		auctionLotId, err := ctx.URLParamInt("auctionLotId")
		if err != nil {
			return err
		}

		resDTO, err := sm.AuctionLotBid.Create(txn, ctx.GetAuth(), uint(auctionLotId), itemDTO)
		if err != nil {
			return err
		}

		ctx.Respond(200, resDTO)

		return nil
	})
	if err != nil {
		ctx.InternalError(err)
	}
}
