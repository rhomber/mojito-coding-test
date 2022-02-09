package errs

import "mojito-coding-test/common/errs"

const (
	//// Auction

	ErrCodeAuctionLotNotFound   errs.ErrorCode = "AUC000"
	ErrCodeAuctionLotInvalid    errs.ErrorCode = "AUC001"
	ErrCodeAuctionLotBidInvalid errs.ErrorCode = "AUC002"
)

//// Auction

var ErrTypeAuctionLotNotFound = errs.NewType(ErrCodeAuctionLotNotFound,
	"Auction Lot Not Found")
var ErrTypeAuctionLotInvalid = errs.NewType(ErrCodeAuctionLotInvalid,
	"Auction Lot Invalid")
var ErrTypeAuctionLotBidInvalid = errs.NewType(ErrCodeAuctionLotBidInvalid,
	"Auction Lot Bid Invalid")
