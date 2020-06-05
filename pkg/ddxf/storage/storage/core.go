package storage

import (
	"encoding/hex"
	"fmt"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	DirData string = "./DataSource/"
)

func UploadDataCore(input *io.StorageUploadInput, ontId string) (output io.StorageUploadOutput) {
	dataSource, err := hex.DecodeString(input.DataSource)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	dataHashOri, err := ddxf.HashObject(input.DataSource)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}
	dataHash := hex.EncodeToString(dataHashOri[:])
	dataPath := DirData + dataHash
	_, err = os.Stat(dataPath)
	if err == nil {
		output.Code = http.StatusInternalServerError
		output.Msg = fmt.Errorf("DataName %s hash %s already exist", input.DataName, dataHash).Error()
		return
	}

	err = ioutil.WriteFile(dataPath, dataSource, 0644)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	dataMap := &io.StorageUploadSave{
		OntId:    ontId,
		DataName: input.DataName,
		DataHash: dataHash,
	}

	err = InsertElt(StorageSaveDataMap, dataMap)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.DataHashUrl = DownloadDataUrl + "/" + dataHash
	output.DataHash = dataHash
	return
}

func DowloadDataCore(input *io.StorageDownloadInput, ontId string) (output io.SorageDownloadOutput) {
	dataPath := DirData + input.DataHash
	_, err := os.Stat(dataPath)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = fmt.Errorf("DataHash %s not exist", input.DataHash).Error()
		return
	}

	dataSource, err := ioutil.ReadFile(dataPath)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.DataSource = hex.EncodeToString(dataSource)
	return
}
