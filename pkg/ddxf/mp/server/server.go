package server

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/middleware/cors"
	"net/http"
)

const (
	addRegistry        = "/ddxf/mp/addRegistry"
	removeRegistry     = "/ddxf/mp/removeRegistry"
	getAuditRule       = "/ddxf/mp/getAuditRule"
	getFee             = "/ddxf/mp/getFee"
	getChallengePeriod = "/ddxf/mp/getChallengePeriod"
	getItemMetaSchema  = "/ddxf/mp/getItemMetaSchema"
	getItemMeta        = "/ddxf/mp/getItemMeta"
	queryItemMetas     = "/ddxf/mp/queryItemMetas"
	PublishItemMeta    = "/ddxf/mp/publishItemMeta"
	Delete             = "/ddxf/mp/delete"
	Update             = "/ddxf/mp/update"
)

func StartMpServer() {
	r := gin.Default()
	r.Use(cors.Cors())
	r.GET(config.Ping, func(context *gin.Context) {
		context.JSON(http.StatusOK, "SUCCESS")
	})
	r.POST(addRegistry, AddRegistryHandler)
	r.POST(removeRegistry, RemoveRegistryHandler)
	r.POST(PublishItemMeta, PublishItemMetaHandler)
	r.POST(getItemMeta, GetItemMetaHandler)
	r.POST(Delete, DeleteHandler)
	r.POST(Update, UpdateHandler)
	r.GET(getAuditRule, GetAuditRuleHandler)
	r.GET(getFee, GetFeeHandler)
	r.GET(getChallengePeriod, GetChallengePeriodHandler)
	r.GET(getItemMetaSchema, GetItemMetaSchemaHandler)
	r.GET(queryItemMetas, QueryItemMetasHandler)
	pri, _ := hex.DecodeString("c19f16785b8f3543bbaf5e1dbb5d398dfa6c85aaad54fc9d71203ce83e505c07")
	MpAccount, _ = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	go r.Run(":" + config.MpPort)
}
