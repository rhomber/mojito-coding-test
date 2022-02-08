package boot

import (
	"math/rand"
	"mojito-coding-test/common/core"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

func Core() func() {
	rand.Seed(time.Now().UnixNano())
	MaxOpenFiles()

	pfCleanup := Profile()

	return func() {
		pfCleanup()
		return
	}
}

func MaxOpenFiles() {
	if runtime.GOOS == "darwin" {
		return
	}

	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		core.Logger.Fatalf("error getting rlimit - %+v", err)
	}
	if rLimit.Cur < rLimit.Max {
		rLimit.Cur = rLimit.Max
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			core.Logger.Fatalf("error setting rlimit - %+v", err)
		} else {
			core.Logger.Debugf("set rlimit to %d", rLimit.Cur)
		}
	}
}

func Profile() func() {
	var debugServer *http.Server
	var cleanupCpuProf func()
	var cleanupTraceProf func()

	if core.Config.GetBool("profile.cpu") {
		f, err := os.Create("cpu.prof")
		if err != nil {
			core.Logger.Fatalf("could not create CPU profile file: %+v", err)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			core.Logger.Fatalf("could not start CPU profile: %+v", err)
		}

		cleanupCpuProf = func() {
			pprof.StopCPUProfile()
			f.Close()
		}
	}

	if core.Config.GetBool("profile.debug.enabled") {
		if core.Config.GetBool("profile.debug.block") {
			runtime.SetBlockProfileRate(1)
		}

		go func() {
			debugAddr := core.Config.GetString("profile.debug.addr")
			debugServer = &http.Server{
				Addr: debugAddr,
			}

			core.Logger.Infof("debug HTTP server: listening on %s", debugAddr)

			if err := debugServer.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					core.Logger.Fatalf("failed to start debug service: %+v", err)
				}
			}
		}()
	}

	return func() {
		if cleanupCpuProf != nil {
			cleanupCpuProf()
		}
		if debugServer != nil {
			core.Logger.Infof("debug HTTP server: shutting down")
			debugServer.Close()
		}
		if cleanupTraceProf != nil {
			cleanupTraceProf()
		}
	}
}
