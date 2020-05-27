package io

// MPEndpointGetFeeInput ...
type MPEndpointGetFeeInput struct {
}

// MPEndpointGetFeeOutput ...
type MPEndpointGetFeeOutput struct {
	Fee
}

// Fee def
type Fee struct {
	ContractAddr string
	Type         TokenType
	Count        uint64
}

// TokenType def
type TokenType byte

const (
	// ONT token
	ONT TokenType = iota
	// ONG token
	ONG
	// OEP4 token
	OEP4
	// OEP5 token
	OEP5
	// OEP8 token
	OEP8
	// OEP68 token
	OEP68
)