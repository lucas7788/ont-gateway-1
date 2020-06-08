package io

import io2 "github.com/zhiqiangxu/ont-gateway/pkg/io"

type StorageCommonServiceOutput struct {
	io2.BaseResp
}

type StorageUploadSave struct {
	OntId    string
	FileName string
}

type StorageUploadOutput struct {
	io2.BaseResp
	DataHashUrl string `json:"dataHashUrl"`
	FileName    string `json:"fileName"`
}

type StorageDownloadInput struct {
	FileName string `json:"fileName"`
}

type SorageDownloadOutput struct {
	io2.BaseResp
	DataSource string `json:"dataSource"`
}
