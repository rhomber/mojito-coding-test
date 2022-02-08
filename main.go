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

	flag.Parse()

	// Core init
	boot.Cfg(*verbose)
	logClose := boot.Logger()
	defer logClose()
	cleanup := boot.Core()
	defer cleanup()

	// App
	a := &app.Application{}

	// Prepare Graph
	core.Populate(a)

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
