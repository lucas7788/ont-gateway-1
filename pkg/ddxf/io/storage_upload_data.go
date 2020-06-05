package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

type StorageCommonServiceOutput struct {
	io2.BaseResp
}

type StorageUploadInput struct {
	DataSource string `json:"dataSource"`
	DataName   string `json:"dataName"`
}

type StorageUploadSave struct {
	OntId    string
	DataName string
	DataHash string
}

type StorageUploadOutput struct {
	io2.BaseResp
	DataHashUrl string `json:"dataHashUrl"`
	DataHash    string `json:"dataHash"`
}

type StorageDownloadInput struct {
	DataHash string `json:"dataHash"`
}

type SorageDownloadOutput struct {
	io2.BaseResp
	DataSource string `json:"dataSource"`
}
