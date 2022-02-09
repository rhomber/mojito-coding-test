package errs

import (
	"mojito-coding-test/common/errs"
	"net/http"
)

////  Auth Namespace

var ErrAuctionLotNotFound = errs.NewForType(ErrTypeAuctionLotNotFound,
	http.StatusBadRequest, "not found")
var ErrAuctionLotInvalid = errs.NewForType(ErrTypeAuctionLotInvalid,
	http.StatusBadRequest, "invalid")
var ErrAuctionLotBidInvalid = errs.NewForType(ErrTypeAuctionLotBidInvalid,
	http.StatusBadRequest, "invalid")
