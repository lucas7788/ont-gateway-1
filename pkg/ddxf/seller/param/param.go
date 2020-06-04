package param

import "github.com/walletsvr/neo/common"

type DataIdInfo struct {
	DataId       string
	DataType     byte
	DataMetaHash common.Uint256
	DataHash     common.Uint256
}
