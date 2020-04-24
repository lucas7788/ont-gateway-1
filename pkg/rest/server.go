package rest

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/intra"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/public"
	"github.com/zhiqiangxu/tableflip"
	"github.com/zhiqiangxu/util/signal"
	"go.uber.org/zap"
)

// Server contains all services
type Server struct {
	IntraApp  *gin.Engine
	PublicApp *gin.Engine
}

// NewServer creates a Server
func NewServer() *Server {

	s := &Server{}

	intraApp := intra.NewApp()
	publicApp := public.NewApp()

	s.IntraApp = intraApp
	s.PublicApp = publicApp

	return s
}

// Start the server
func (s *Server) Start() {

	conf := config.Load()

	upg, err := tableflip.New(tableflip.Options{PIDFile: conf.RestConfig.PIDFile})
	if err != nil {
		instance.Logger().Error("tableflip", zap.Error(err))
		return
	}
	defer upg.Stop()

	intraServer := &http.Server{Addr: conf.RestConfig.IntraAddr, Handler: s.IntraApp, ReadTimeout: conf.RestConfig.ReadTimeout}
	publicServer := &http.Server{Addr: conf.RestConfig.PublicAddr, Handler: s.PublicApp, ReadTimeout: conf.RestConfig.ReadTimeout}

	lnIntra, err := upg.Fds.Listen("tcp", conf.RestConfig.IntraAddr)
	if err != nil {
		instance.Logger().Error("Listen IntraAddr", zap.Error(err))
		return
	}
	lnPublic, err := upg.Fds.Listen("tcp", conf.RestConfig.PublicAddr)
	if err != nil {
		instance.Logger().Error("Listen PublicAddr", zap.Error(err))
		return
	}

	var g run.Group
	g.Add(func() error {
		return intraServer.Serve(lnIntra)
	}, func(error) {
		intraServer.Shutdown(context.Background())
	})
	g.Add(func() error {
		return publicServer.Serve(lnPublic)
	}, func(error) {
		publicServer.Shutdown(context.Background())
	})

	groupDoneCh := make(chan struct{})
	go func() {
		err := g.Run()
		if err != nil {
			instance.Logger().Error("Group.Run", zap.Error(err))
		}
		close(groupDoneCh)
	}()

	if err := upg.Ready(); err != nil {
		instance.Logger().Error("upg.Ready", zap.Error(err))
		return
	}
	instance.Logger().Error("child ready to serve")

	shutdownServers := func() (err error) {
		closeTimeout := time.Second * 2
		{
			ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
			defer cancel()
			err = intraServer.Shutdown(ctx)
			if err != nil {
				instance.Logger().Error("shutdownServers intraServer.Shutdown", zap.Error(err))
			}
		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
			defer cancel()
			err2 := publicServer.Shutdown(ctx)
			if err2 != nil {
				instance.Logger().Error("shutdownServers publicServer.Shutdown", zap.Error(err2))
				err = err2
			}
		}
		return
	}

	signal.SetupHandler(func(s os.Signal) {
		// upgrade on signal
		instance.Logger().Error("on upgrade", zap.Bool("GracefulUpgrade", config.Load().RestConfig.GracefulUpgrade))

		if !config.Load().RestConfig.GracefulUpgrade {
			err := shutdownServers()
			if err != nil {
				instance.Logger().Error("shutdownServers", zap.Error(err))
			}
			os.Exit(1)
		}

		instance.Logger().Error("upgrade start")
		err := upg.Upgrade()
		instance.Logger().Error("upgrade end", zap.Error(err))
	}, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGQUIT)

	// ready to exit
	select {
	case <-upg.Exit():
	case <-groupDoneCh:
	}
	instance.Logger().Error("parent prepare for exit")

	err = shutdownServers()
	if err != nil {
		instance.Logger().Error("shutdownServers err", zap.Error(err))
	}

	instance.Logger().Error("parent quit ok")

}
