package boot

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"strings"
)

const dbMigrationPath string = "./db/migrations"

func doMigration(logger *logrus.Entry, dbConn *sql.DB, rollback bool) error {
	var err error

	if rollback {
		logger.Infof("Rolling back %s", dbMigrationPath)
		err = goose.Down(dbConn, dbMigrationPath)
	} else {
		logger.Infof("Applying %s", dbMigrationPath)
		err = goose.Up(dbConn, dbMigrationPath)
	}
	if err != nil {
		// Looks horrible and hard to read otherwise.
		fmt.Printf("[ERROR OUTPUT] failed to perform migration\n\n%+v",
			strings.ReplaceAll(err.Error(), "\\n", "\n"))

		return fmt.Errorf("failed to perform migration: %+v", err)
	}

	return err
}
