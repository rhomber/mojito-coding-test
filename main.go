package main

import (
	"flag"
	"mojito-coding-test/app"
	"mojito-coding-test/common/boot"
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/core"
)

func main() {
	// Parse switches
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	rollback := flag.Bool("rollback-migration", false, "Rollback last migration")

	flag.Parse()

	// Core init
	boot.Cfg(*verbose)
	logClose := boot.Logger()
	defer logClose()
	cleanup := boot.Core()
	defer cleanup()

	// Db init
	db := boot.Sqlite3(*rollback)

	// App
	a := &app.Application{}

	// Prepare Graph
	core.Populate(a, db)

	// Prepare App
	handlers, err := a.Init()
	if err != nil {
		core.Logger.Fatal(err)
	}

	// Start HTTP
	chttp.NewServer(core.Logger, core.Config, handlers,
		// OnStart
		func() {
		},
		// OnShutdown
		func(isRestart bool) {
		},
	).Start()
}
