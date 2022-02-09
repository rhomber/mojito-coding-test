package handler

import (
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/service"
	"mojito-coding-test/common/chttp"
	"net/http"
)

func GetAuctionLots(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	err := ctx.GetDb().Transaction(func(txn *gorm.DB) error {
		sm := ctx.GetServiceManager().(*service.Manager)

		lotDTOs, err := sm.AuctionLot.List(txn)
		if err != nil {
			return err
		}

		ctx.Respond(200, lotDTOs)

		return nil
	})
	if err != nil {
		ctx.InternalError(err)
	}
}

func PostAuctionLot(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	var itemDTO dto.CreateAuctionLot
	if !ctx.MustDecode(&itemDTO) {
		return
	}

	err := ctx.GetDb().Transaction(func(txn *gorm.DB) error {
		sm := ctx.GetServiceManager().(*service.Manager)

		resDTO, err := sm.AuctionLot.Create(txn, itemDTO)
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
