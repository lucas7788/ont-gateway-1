package model

// ResourceVersion model
type ResourceVersion struct {
	App   int    `bson:"app" json:"app"`
	ID    string `bson:"id" json:"id"`
	Block uint32 `bson:"block" json:"block"`
	Desc  string `bson:"desc" json:"desc"`
	Hash  string `bson:"hash" json:"hash"`
}
