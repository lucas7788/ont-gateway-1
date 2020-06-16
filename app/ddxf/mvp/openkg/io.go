package main

// PublishInput ...
type PublishInput struct {
	ReqID    string
	OpenKGID string
	UserID   string // seller id
	Item     map[string]interface{}
	Datas    []map[string]interface{}
	Delete   bool
}

// PublishOutput ...
type PublishOutput struct {
	Code  int
	Msg   string
	ReqID string
}

// BuyAndUseInput ...
type BuyAndUseInput struct {
	ReqID    string
	OpenKGID string
	DataID   string
	UserID   string // buyer id
}

// BuyAndUseOutput ...
type BuyAndUseOutput struct {
	Code  int
	Msg   string
	ReqID string
}
