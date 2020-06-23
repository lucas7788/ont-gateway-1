package server

type GenerateOntIdInput struct {
	ReqID  string `bson:"req_id" json:"req_id"`
	UserId string `bson:"user_id" json:"user_id"`
}

type GenerateOntIdOutput struct {
	Code  int    `bson:"code" json:"code"`
	Msg   string `bson:"msg" json:"msg"`
	ReqID string `bson:"req_id" json:"req_id"`
	OntId string `bson:"ont_id" json:"ont_id"`
}

// PublishInput ...
type PublishInput struct {
	ReqID     string                   `bson:"req_id" json:"req_id"`
	OpenKGID  string                   `bson:"openkg_id" json:"openkg_id"`
	UserID    string                   `bson:"user_id" json:"user_id"` // seller id
	Item      map[string]interface{}   `bson:"item" json:"item"`
	Datas     []map[string]interface{} `bson:"datas" json:"datas"`
	Delete    bool                     `bson:"delete" json:"delete"`
	OnChainId string                   `bson:"on_chain_id" json:"on_chain_id"`
}

// PublishOutput ...
type PublishOutput struct {
	Code  int    `bson:"code" json:"code"`
	Msg   string `bson:"msg" json:"msg"`
	ReqID string `bson:"req_id" json:"req_id"`
}

// BuyAndUseInput ...
type BuyAndUseInput struct {
	ReqID    string `bson:"req_id" json:"req_id"`
	OpenKGID string `bson:"openkg_id" json:"openkg_id"`
	DataID   string `bson:"data_id" json:"data_id"`
	UserID   string `bson:"user_id" json:"user_id"` // buyer id
}

// BuyAndUseOutput ...
type BuyAndUseOutput struct {
	Code  int    `bson:"code" json:"code"`
	Msg   string `bson:"msg" json:"msg"`
	ReqID string `bson:"req_id" json:"req_id"`
}

// BuyAndUseInput ...
type DeleteInput struct {
	ReqID      string `bson:"req_id" json:"req_id"`
	OpenKGID   string `bson:"openkg_id" json:"openkg_id"`
	ResourceID string `bson:"data_id" json:"data_id"`
	UserID     string `bson:"user_id" json:"user_id"` // buyer id
}

// BuyAndUseInput ...
type DeleteOutput struct {
	ReqID string `bson:"req_id" json:"req_id"`
	Code  int    `bson:"code" json:"code"`
	Msg   string `bson:"msg" json:"msg"`
}
