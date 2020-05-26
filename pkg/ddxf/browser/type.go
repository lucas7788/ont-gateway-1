package browser

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Browser ...
type Browser interface {
	EditDataMeta(io.BrowserEditDataMetaInput) io.BrowserEditDataMetaOutput
	EditTokenMeta(io.BrowserEditTokenMetaInput) io.BrowserEditTokenMetaOutput
	VerifyDataAndToken(io.BrowserVerifyDataAndTokenInput) io.BrowserVerifyDataAndTokenOutput

	ConstructBuyDtokenTx(io.BrowserConstructBuyDtokenTxInput) io.BrowserConstructBuyDtokenTxOutput
}
