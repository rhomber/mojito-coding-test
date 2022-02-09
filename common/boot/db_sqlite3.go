package boot

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"mojito-coding-test/common/core"
	"os"
)

func Sqlite3(rollback bool) *gorm.DB {
	if err := migrateSqlite3Db(rollback); err != nil {
		core.Logger.Fatalf("failed to init db - %+v", err)
	}

	if rollback {
		// Exit after rollback.
		os.Exit(0)
	}

	// Init Gorm
	db, err := gorm.Open(sqlite.Open(core.Config.GetDbFile()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		core.Logger.Fatalf("failed to connect to db - %+v", err)
	}

	return db
}

func migrateSqlite3Db(rollback bool) error {
	cfg := core.Config
	logger := core.Logger

	dbFile := core.Config.GetDbFile()

	if !rollback {
		logger.Infof("Performing migrations for '%s' ...", dbFile)
	} else {
		logger.Infof("Performing rollback for '%s' ...", dbFile)
	}

	dbConn, err := sql.Open("sqlite3", cfg.GetDsn())
	defer dbConn.Close()

	if err != nil {
		return fmt.Errorf("failed to connect to db ('%s') for migration: %v", dbFile, err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		return errors.Wrap(err, "setting dialect")
	}

	if err := doMigration(logger, dbConn, rollback); err != nil {
		return err
	}

	return nil
}
