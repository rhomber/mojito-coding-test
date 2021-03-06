package handler

import (
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/service"
	"mojito-coding-test/common/chttp"
	"net/http"
)

func GetUsers(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	err := ctx.GetDb().Transaction(func(txn *gorm.DB) error {
		sm := ctx.GetServiceManager().(*service.Manager)

		userDTOs, err := sm.User.List(txn)
		if err != nil {
			return err
		}

		ctx.Respond(200, userDTOs)

		return nil
	})
	if err != nil {
		ctx.InternalError(err)
	}
}

func PostUser(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	var itemDTO dto.CreateUser
	if !ctx.MustDecode(&itemDTO) {
		return
	}

	err := ctx.GetDb().Transaction(func(txn *gorm.DB) error {
		sm := ctx.GetServiceManager().(*service.Manager)

		resDTO, err := sm.User.Create(txn, itemDTO)
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
