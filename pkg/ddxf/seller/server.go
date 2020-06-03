package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/restful"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/sellerconfig"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/seller/service"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

const (
	SaveDataMetaUrl            = "/ddxf/seller/saveDataMeta"
	SaveTokenMetaUrl           = "/ddxf/seller/saveTokenMeta"
	PublishMPItemMetaUrl       = "/ddxf/seller/publishMPItemMeta"
	UseTokenUrl                = "/ddxf/seller/useToken"
	PublishItemMetaUrl         = "/ddxf/seller/publishItemMeta"
	getQrCodeDataByQrCodeIdUrl = "ddxf/seller/getQrCodeDataByQrCodeId"
	qrCodeCallbackSendTxUrl    = "ddxf/seller/qrCodeCallbackSendTx"
)

var (
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: uint(log.InfoLog),
	}
	ConfigfileFlag = cli.StringFlag{
		Name:   "config",
		Usage:  "specify configfile",
		Value:  "config.json",
		EnvVar: "CONFIG_FILE",
	}
)

func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}

func startDDXFSeller(ctx *cli.Context) error {
	// init log level
	logLevel := ctx.GlobalInt(GetFlagName(LogLevelFlag))
	log.InitLog(logLevel, log.Stdout)

	// init config
	configName := ctx.String(GetFlagName(ConfigfileFlag))
	sellerConfig, err := sellerconfig.InitSellerConfig(configName)
	if err != nil {
		return err
	}
	*sellerconfig.DefSellerConfig = *sellerConfig
	service.InitSellerImpl()
	StartSellerServer()
	waitToExit()
	return nil
}

func StartSellerServer() {
	r := gin.Default()
	//r.Use(middleware.JWT)
	r.POST(SaveDataMetaUrl, restful.SaveDataMetaHandle)
	r.POST(SaveTokenMetaUrl, restful.SaveTokenMetaHandle)
	r.POST(PublishMPItemMetaUrl, restful.PublishMPItemMetaHandle)
	r.POST(UseTokenUrl, restful.UseTokenHandler)
	r.POST(PublishItemMetaUrl, restful.PublishMetaHandler)
	r.POST(getQrCodeDataByQrCodeIdUrl, restful.GetQrCodeDataByQrCodeId)
	r.POST(qrCodeCallbackSendTxUrl, restful.GrCodeCallbackSendTx)
	go r.Run(":" + config.SellerPort)
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("ddxf seller server received exit signal: %s.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "ddxf seller CLI"
	app.Action = startDDXFSeller
	app.Version = "1.0"
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		LogLevelFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
