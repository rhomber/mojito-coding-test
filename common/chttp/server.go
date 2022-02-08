package chttp

import (
	"context"
	"github.com/sirupsen/logrus"
	"mojito-coding-test/common/config"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

func NewServer(logger *logrus.Entry, cfg *config.Config,
	handlers http.Handler, onStart func(), onShutdown func(bool)) *Server {
	return &Server{
		Logger:     logger,
		Config:     cfg,
		Handlers:   handlers,
		OnStart:    onStart,
		OnShutdown: onShutdown,
	}
}

type Server struct {
	Logger     *logrus.Entry
	Config     *config.Config
	Handlers   http.Handler
	OnStart    func()
	OnShutdown func(bool)

	srv            *http.Server
	isShuttingDown bool
}

func (s *Server) Restart() {
	s.Shutdown(true)
	s.Start()
}

func (s *Server) Shutdown(isRestart bool) {
	if s.srv == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			s.Logger.Errorf("recovered while shutting down, ignoring - %+v\n%s",
				r, string(debug.Stack()))
		}
	}()

	s.isShuttingDown = true

	shutdownWait := s.Config.GetDuration("http.shutdown.waitTimeOut")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownWait)
	defer cancel()
	s.srv.Shutdown(ctx)
	if s.OnShutdown != nil {
		s.OnShutdown(isRestart)
	}
}

func (s *Server) Start() {
	defer func() {
		if r := recover(); r != nil {
			s.Logger.Errorf("HTTP server recovered (will restart): %+v", r)
			time.Sleep(time.Second * 2)
			s.Restart()
		}
	}()

	s.isShuttingDown = false

	if s.OnStart != nil {
		s.OnStart()
	}

	serverAddress := s.Config.GetString("http.addr")

	port := os.Getenv("PORT")
	if port != "" {
		serverAddress = ":" + port
	}

	s.srv = &http.Server{Addr: serverAddress, Handler: s.Handlers}

	sigWg := sync.WaitGroup{}
	sigWg.Add(1)

	sigQuitCh := make(chan bool)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		select {
		case <-sigQuitCh:
			break
		case sig := <-gracefulStop:
			s.Logger.Infof("caught sig: %+v, shutting down!", sig)
			s.Shutdown(false)
		}

		sigWg.Done()
	}()

	s.Logger.Infof("HTTP server: listening on %s", serverAddress)

	if err := s.srv.ListenAndServe(); err != nil {
		if !s.isShuttingDown {
			s.Logger.Fatal(err)
		}
	}

	close(sigQuitCh)
	sigWg.Wait()
}
