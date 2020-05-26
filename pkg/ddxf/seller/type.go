package seller

import "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"

// Seller ...
type Seller interface {
	SaveDataMeta(io.SellerSaveDataMetaInput) io.SellerSaveDataMetaOutput
	SaveTokenMeta(io.SellerSaveTokenMetaInput) io.SellerSaveTokenMetaOutput
	PublishMPItemMeta(io.SellerPublishMPItemMetaInput) io.SellerPublishMPItemMetaOutput
}
