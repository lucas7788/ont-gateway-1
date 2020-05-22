package io

// JSONLDAlignInput ...
type JSONLDAlignInput struct {
	A map[string]interface{}
	B map[string]interface{}
}

// JSONLDAlignOutput ...
type JSONLDAlignOutput struct {
	BaseResp
	A2B map[string]interface{}
}
