package intra

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	g "github.com/zhiqiangxu/ont-gateway/pkg/rest/gin"
	"github.com/zhiqiangxu/ont-gateway/pkg/rest/middleware"
)

// NewApp for intra app
func NewApp() *gin.Engine {
	r := g.New()

	pprof.Register(r)
	// metrics内网展示
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// this api is for test purpose
	aksk := r.Group("/aksk")
	aksk.Use(middleware.AkSk)
	aksk.POST("/hello", AkSkHello)

	addon := aksk.Group("/addon")

	// per tenant stuff

	addon.POST("/config", UpsertAddonConfig)
	addon.GET("/config", GetAddonConfig)
	addon.DELETE("/config", DeleteAddonConfig)

	addon.POST("/deploy", DeployAddon)
	addon.GET("/deploy/check", CheckDeploy)
	addon.POST("/shell", Shell)

	tx := aksk.Group("/tx")
	tx.POST("/poll", EnqueTx)
	tx.DELETE("/poll", DequeTx)

	payment := aksk.Group("/payment")
	payment.POST("/config", CreatePaymentConfig)
	payment.DELETE("/config", DeletePaymentConfig)
	payment.POST("/order", CreatePaymentOrder)

	wallet := aksk.Group("/wallet")
	wallet.GET("/", GetWallet)

	resourceVersion := aksk.Group("/rv")
	resourceVersion.POST("/", UpdateResource)
	resourceVersion.GET("/", GetResource)

	return r
}

func sendoutput(c *gin.Context, code int, output interface{}) {
	if code == 0 {
		code = http.StatusOK
	}

	c.JSON(code, output)
}
