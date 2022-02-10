package service

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/model"
	"mojito-coding-test/app/data/types"
	"mojito-coding-test/common/config"
	cdto "mojito-coding-test/common/data/dto"
	"os"
	"testing"
	"time"
)

const (
	bidIncrement = 100
)

type AuctionLotBidSuite struct {
	suite.Suite

	db     *gorm.DB
	dbMock sqlmock.Sqlmock
	svc    *AuctionLotBid
}

// Admittedly this would be better with a mocked repository vs. mocked SQL (or in addition to).
// I am also aware of how a few of the 'unit' tests bleed into other unrelated code (using interfaces and mocking
// would be ideal really).

func (s *AuctionLotBidSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.dbMock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(s.T(), err)

	// GORM expected queries.
	s.dbMock.ExpectQuery("select sqlite_version()").
		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).
			AddRow(3))

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Silent,
			//LogLevel: logger.Info,
			Colorful: true,
		},
	)

	s.db, err = gorm.Open(sqlite.Dialector{
		Conn: db,
	}, &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	require.NoError(s.T(), err)

	cfg := config.Config{Viper: viper.New()}
	cfg.Set("bid.increment", bidIncrement)

	lgr := logrus.NewEntry(logrus.New())

	s.svc = &AuctionLotBid{
		Logger:    lgr,
		Validator: validator.New(),
		Config:    &cfg,
		Db:        s.db,
		AuctionLot: &AuctionLot{
			Config:    &cfg,
			Validator: validator.New(),
			Db:        s.db,
		},
	}
}

func (s *AuctionLotBidSuite) TestFindCurrentBid() {
	now := time.Now()

	tables := []struct {
		id           uint
		auctionLotId uint
		userId       uint
		bid          uint
		bidType      types.BidType
		createdAt    time.Time
		expError     string
	}{
		{id: 2, auctionLotId: 10, userId: 12, bid: 200, bidType: types.BidTypeUser, createdAt: now},
		{id: 3, auctionLotId: 10, userId: 14, bid: 900, bidType: types.BidTypeMaxBid, createdAt: now},
		{auctionLotId: 11, expError: "record not found"},
	}

	for _, tst := range tables {
		expQ := s.dbMock.ExpectQuery("SELECT * FROM `auction_lot_bid` " +
			"WHERE auction_lot_id=? AND `auction_lot_bid`.`deleted_at` IS NULL " +
			"ORDER BY bid DESC,`auction_lot_bid`.`id` LIMIT 1").
			WithArgs(tst.auctionLotId)

		rowsMock := sqlmock.NewRows([]string{"id", "auction_lot_id", "user_id", "type", "bid", "created_at"})

		if tst.id > 0 {
			expQ.WillReturnRows(rowsMock.
				AddRow(tst.id, tst.auctionLotId, tst.userId, string(tst.bidType), tst.bid, tst.createdAt))
		} else {
			expQ.WillReturnRows(rowsMock)
		}

		res, err := s.svc.FindCurrentBid(s.db, tst.auctionLotId)

		require.NoError(s.T(), s.dbMock.ExpectationsWereMet())

		if tst.expError != "" {
			require.Nil(s.T(), deep.Equal(model.AuctionLotBid{}, res))
			require.Equal(s.T(), tst.expError, err.Error())
		} else {
			require.NoError(s.T(), err)
			require.Nil(s.T(), deep.Equal(model.AuctionLotBid{
				Model: gorm.Model{
					ID:        tst.id,
					CreatedAt: tst.createdAt,
				},
				AuctionLotId: tst.auctionLotId,
				UserId:       tst.userId,
				Type:         tst.bidType,
				Bid:          tst.bid,
			}, res))
		}
	}
}

func (s *AuctionLotBidSuite) TestFindMaxBid() {
	now := time.Now()

	tables := []struct {
		id           uint
		auctionLotId uint
		userId       uint
		maxBid       uint
		active       bool
		createdAt    time.Time
		expError     string
	}{
		{id: 2, auctionLotId: 10, userId: 12, maxBid: 200, active: true, createdAt: now},
		{id: 3, auctionLotId: 10, userId: 14, maxBid: 900, active: true, createdAt: now},
		{auctionLotId: 11, expError: "record not found"},
	}

	for _, tst := range tables {
		expQ := s.dbMock.ExpectQuery("SELECT * FROM `auction_lot_bid_max` "+
			"WHERE (auction_lot_id=? AND user_id=? AND active=true) "+
			"AND `auction_lot_bid_max`.`deleted_at` IS NULL ORDER BY `auction_lot_bid_max`.`id` LIMIT 1").
			WithArgs(tst.auctionLotId, tst.userId)

		rowsMock := sqlmock.NewRows([]string{"id", "auction_lot_id", "user_id", "max_bid", "active", "created_at"})

		if tst.id > 0 {
			expQ.WillReturnRows(rowsMock.
				AddRow(tst.id, tst.auctionLotId, tst.userId, tst.maxBid, tst.active, tst.createdAt))
		} else {
			expQ.WillReturnRows(rowsMock)
		}

		res, err := s.svc.FindMaxBid(s.db, tst.auctionLotId, tst.userId)

		require.NoError(s.T(), s.dbMock.ExpectationsWereMet())

		if tst.expError != "" {
			require.Nil(s.T(), deep.Equal(model.AuctionLotBidMax{}, res))
			require.Equal(s.T(), tst.expError, err.Error())
		} else {
			require.NoError(s.T(), err)
			require.Nil(s.T(), deep.Equal(model.AuctionLotBidMax{
				Model: gorm.Model{
					ID:        tst.id,
					CreatedAt: tst.createdAt,
				},
				AuctionLotId: tst.auctionLotId,
				UserId:       tst.userId,
				MaxBid:       tst.maxBid,
				Active:       tst.active,
			}, res))
		}
	}
}

/** Skipped due to time constraints.

func (s *AuctionLotBidSuite) TestListBids() {

}

func (s *AuctionLotBidSuite) TestListMaxBids() {

}

func (s *AuctionLotBidSuite) TestList() {

}
*/

func (s *AuctionLotBidSuite) TestCreate() {
	var now = time.Now()
	var aStartDate = time.Now().Add(-time.Hour * 24)
	var aEndDate = time.Now().Add(time.Hour * 24)

	tables := []struct {
		authUserId          uint
		auctionLotId        uint
		auctionLotStartDate time.Time
		auctionLotEndDate   time.Time
		createBid           uint
		createMaxBid        uint
		curBidId            uint
		curBidUserId        uint
		curBid              uint
		curBidType          types.BidType
		curMaxBidId         uint
		curMaxBid           uint
		expNoQueryLot       bool
		expNoQueryBid       bool
		expSuccess          bool
		expInsertBid        bool
		expInsertMaxBid     bool
		expBid              uint
		expBidType          types.BidType
		expBidUserId        uint
		expMaxBid           uint
		expError            string
	}{
		// Validation
		{authUserId: 10, auctionLotId: 2, expError: "validation failed for creation of new auction lot bid: " +
			"Key: 'CreateAuctionLotBid.Bid' Error:Field validation for 'Bid' failed on the 'required' tag",
			expNoQueryLot: true, expNoQueryBid: true},
		{authUserId: 10, auctionLotId: 11, createBid: 12, createMaxBid: 100,
			expError:      fmt.Sprintf("[AUC002] Auction Lot Bid Invalid: bid must be increment of %d", bidIncrement),
			expNoQueryLot: true, expNoQueryBid: true},
		{authUserId: 10, auctionLotId: 11, createBid: 200, createMaxBid: 110,
			expError:      fmt.Sprintf("[AUC002] Auction Lot Bid Invalid: max bid must be increment of %d", bidIncrement),
			expNoQueryLot: true, expNoQueryBid: true},
		{authUserId: 10, auctionLotId: 11, createBid: 100, createMaxBid: 300, auctionLotStartDate: time.Now().Add(time.Hour * 24),
			expError:      "[AUC001] Auction Lot Invalid: start time is in the future",
			expNoQueryBid: true},
		{authUserId: 10, auctionLotId: 11, createBid: 100, createMaxBid: 300, auctionLotStartDate: time.Now().Add(time.Hour * 24),
			expError:      "[AUC001] Auction Lot Invalid: start time is in the future",
			expNoQueryBid: true},
		{authUserId: 10, auctionLotId: 11, createBid: 100, createMaxBid: 300, auctionLotStartDate: aStartDate,
			auctionLotEndDate: time.Now().Add(-time.Hour * 24), expError: "[AUC001] Auction Lot Invalid: end time is in the past",
			expNoQueryBid: true},

		// Scenarios
		// 0: No Current Bids
		{authUserId: 10, auctionLotId: 11, createBid: 100, createMaxBid: 300, auctionLotStartDate: aStartDate,
			auctionLotEndDate: aEndDate, expSuccess: true, expBid: 100, expBidType: types.BidTypeUser, expBidUserId: 10,
			expMaxBid: 300, expInsertBid: true, expInsertMaxBid: true},
		// 1: Update max bids only
		{authUserId: 10, auctionLotId: 11, createBid: 100, createMaxBid: 300, auctionLotStartDate: aStartDate,
			auctionLotEndDate: aEndDate, curBidId: 9, curBidUserId: 10, curBid: 100, curBidType: types.BidTypeUser,
			curMaxBidId: 9, curMaxBid: 200, expSuccess: true, expBid: 100, expBidType: types.BidTypeUser,
			expBidUserId: 10, expMaxBid: 300, expInsertBid: false, expInsertMaxBid: true},
		// 2: Bid is greater than current max bid
		{authUserId: 10, auctionLotId: 11, createBid: 300, createMaxBid: 400, auctionLotStartDate: aStartDate,
			auctionLotEndDate: aEndDate, curBidId: 9, curBidUserId: 20, curBid: 200, curBidType: types.BidTypeUser,
			curMaxBidId: 9, curMaxBid: 200, expSuccess: true, expBid: 300, expBidType: types.BidTypeUser,
			expBidUserId: 10, expMaxBid: 400, expInsertBid: true, expInsertMaxBid: true},
		// 3: Current max bid exceeds submitted max bid.
		{authUserId: 10, auctionLotId: 11, createBid: 300, createMaxBid: 400, auctionLotStartDate: aStartDate,
			auctionLotEndDate: aEndDate, curBidId: 9, curBidUserId: 20, curBid: 200, curBidType: types.BidTypeUser,
			curMaxBidId: 9, curMaxBid: 900, expSuccess: false, expBid: 400, expBidType: types.BidTypeMaxBid,
			expBidUserId: 20, expMaxBid: 900, expInsertBid: true, expInsertMaxBid: false},
		// 4: Submitted max exceeds current max
		{authUserId: 10, auctionLotId: 11, createBid: 100, createMaxBid: 300, auctionLotStartDate: aStartDate,
			auctionLotEndDate: aEndDate, curBidId: 9, curBidUserId: 20, curBid: 200, curBidType: types.BidTypeUser,
			curMaxBidId: 9, curMaxBid: 200, expSuccess: true, expBid: 300, expBidType: types.BidTypeMaxBid,
			expBidUserId: 10, expMaxBid: 300, expInsertBid: true, expInsertMaxBid: true},
		// 5: Submitted isn't greater in bid nor max bid
		{authUserId: 10, auctionLotId: 11, createBid: 300, createMaxBid: 400, auctionLotStartDate: aStartDate,
			auctionLotEndDate: aEndDate, curBidId: 9, curBidUserId: 20, curBid: 400, curBidType: types.BidTypeUser,
			curMaxBidId: 9, curMaxBid: 900, expSuccess: false, expBid: 400, expBidType: types.BidTypeUser,
			expBidUserId: 20, expMaxBid: 900, expInsertBid: false, expInsertMaxBid: false},
	}

	for _, tst := range tables {
		// Setup
		auth := cdto.Auth{UserId: tst.authUserId}
		createDTO := dto.CreateAuctionLotBid{
			Bid:    tst.createBid,
			MaxBid: tst.createMaxBid,
		}
		expBidId := tst.curBidId
		expMaxBidId := tst.curMaxBidId

		// Reset Mocks
		s.SetupTest()

		if !tst.expNoQueryLot {
			// Auction Lot
			s.dbMock.ExpectQuery("SELECT * FROM `auction_lot` " +
				"WHERE `auction_lot`.`id` = ? AND `auction_lot`.`deleted_at` IS NULL " +
				"ORDER BY `auction_lot`.`id` LIMIT 1").
				WithArgs(tst.auctionLotId).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "start_time", "end_time"}).
					AddRow(tst.auctionLotId, "Test", tst.auctionLotStartDate, tst.auctionLotEndDate))
		}

		if !tst.expNoQueryBid {
			// Find Current Bid
			mbExpQ := s.dbMock.ExpectQuery("SELECT * FROM `auction_lot_bid` " +
				"WHERE auction_lot_id=? AND `auction_lot_bid`.`deleted_at` IS NULL " +
				"ORDER BY bid DESC,`auction_lot_bid`.`id` LIMIT 1").
				WithArgs(tst.auctionLotId)

			rowsMock := sqlmock.NewRows([]string{"id", "auction_lot_id", "user_id", "type", "bid", "created_at"})

			if tst.curBidId > 0 {
				mbExpQ.WillReturnRows(rowsMock.
					AddRow(tst.curBidId, tst.auctionLotId, tst.curBidUserId, string(tst.curBidType), tst.curBid,
						now))
			} else {
				mbExpQ.WillReturnRows(rowsMock)
			}
		}

		// Find Max Bid
		if tst.curBidId > 0 {
			s.dbMock.ExpectQuery("SELECT * FROM `auction_lot_bid_max` "+
				"WHERE (auction_lot_id=? AND user_id=? AND active=true) "+
				"AND `auction_lot_bid_max`.`deleted_at` IS NULL ORDER BY `auction_lot_bid_max`.`id` LIMIT 1").
				WithArgs(tst.auctionLotId, tst.curBidUserId).
				WillReturnRows(sqlmock.NewRows([]string{"id", "auction_lot_id", "user_id", "max_bid", "active"}).
					AddRow(tst.curMaxBidId, tst.auctionLotId, tst.authUserId, tst.curMaxBid, true))
		}

		// Insert Bid & Max Bid
		if tst.expInsertBid || tst.expInsertMaxBid {
			if tst.expInsertBid {
				expBidId++

				s.dbMock.ExpectBegin()
				s.dbMock.ExpectExec("INSERT INTO `auction_lot_bid` (`created_at`,`updated_at`,`deleted_at`,"+
					"`auction_lot_id`,`user_id`,`type`,`bid`) VALUES (?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tst.auctionLotId, tst.expBidUserId,
						tst.expBidType, tst.expBid).
					WillReturnResult(sqlmock.NewResult(int64(expBidId), 1))
				s.dbMock.ExpectCommit()
			}

			if tst.expInsertMaxBid {
				expMaxBidId++

				// Reset
				s.dbMock.ExpectExec("UPDATE auction_lot_bid_max SET active=NULL "+
					"WHERE auction_lot_id=? AND user_id=?").
					WithArgs(tst.auctionLotId, tst.expBidUserId).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Insert
				s.dbMock.ExpectBegin()
				s.dbMock.ExpectExec("INSERT INTO `auction_lot_bid_max` (`created_at`,`updated_at`,"+
					"`deleted_at`,`auction_lot_id`,`user_id`,`max_bid`,`active`) VALUES (?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tst.auctionLotId, tst.expBidUserId,
						tst.expMaxBid, true).
					WillReturnResult(sqlmock.NewResult(int64(expMaxBidId), 1))
				s.dbMock.ExpectCommit()
			}
		}

		res, err := s.svc.Create(s.db, auth, tst.auctionLotId, createDTO)

		require.NoError(s.T(), s.dbMock.ExpectationsWereMet())

		if tst.expError != "" {
			require.Nil(s.T(), deep.Equal(dto.CreateAuctionLotBidResult{}, res))
			require.Equal(s.T(), tst.expError, err.Error())
		} else {
			require.NoError(s.T(), err)
			require.Nil(s.T(), deep.Equal(dto.CreateAuctionLotBidResult{
				Success:      tst.expSuccess,
				Id:           expBidId,
				AuctionLotId: tst.auctionLotId,
				UserId:       tst.expBidUserId,
				Type:         tst.expBidType,
				Bid:          tst.expBid,
				MaxBid:       tst.expMaxBid,
				CreatedAt:    res.CreatedAt, // Hard to test as GORM creates this.
			}, res))
		}
	}

}

func TestAuctionLotBidSuite(t *testing.T) {

	suite.Run(t, new(AuctionLotBidSuite))
}
