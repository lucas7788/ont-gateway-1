package storage

import (
	"encoding/hex"
	"fmt"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/common"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	io2 "io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	DirData           string = "./DataSource"
	StorageFilePrefix string = "StorageFilePrefix"
)

func UploadDataCore(file multipart.File, ontId string) (output io.StorageUploadOutput) {
	fileName := common.GenerateUUId(StorageFilePrefix)
	_, err := os.Stat(DirData)
	if err != nil {
		err = os.Mkdir(DirData, os.ModePerm)
		if err != nil {
			output.Code = http.StatusBadRequest
			output.Msg = err.Error()
			return
		}
	}
	dataPath := fmt.Sprintf("%s%s%s", DirData, string(os.PathSeparator), fileName)
	out, err := os.Create(dataPath)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	defer out.Close()
	_, err = io2.Copy(out, file)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	dataMap := &io.StorageUploadSave{
		OntId:    TestOntId,
		FileName: fileName,
	}

	err = InsertElt(StorageSaveDataMap, dataMap)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = err.Error()
		return
	}

	output.DataHashUrl = DownloadDataPrefix + fileName
	output.FileName = fileName

	return
}

func DowloadDataCore(input *io.StorageDownloadInput, ontId string) (output io.SorageDownloadOutput) {
	dataPath := fmt.Sprintf("%s%s%s", DirData, string(os.PathSeparator), input.FileName)
	_, err := os.Stat(dataPath)
	if err != nil {
		output.Code = http.StatusInternalServerError
		output.Msg = fmt.Errorf("DataHash %s not exist", input.FileName).Error()
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
