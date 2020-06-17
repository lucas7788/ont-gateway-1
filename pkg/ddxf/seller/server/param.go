package server

type ItemMeta struct {
	ItemMetaHash string                 `bson:"itemMetaHash" json:"itemMetaHash"`
	ItemMetaData map[string]interface{} `bson:"itemMetaData" json:"itemMetaData"`
}
