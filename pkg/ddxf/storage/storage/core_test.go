package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"testing"
)

func TestUploadDataCore(t *testing.T) {
	upinput := io.StorageUploadInput{
		DataSource: "0102030405",
		DataName:   "filetest",
	}
	output := UploadDataCore(&upinput, TestOntId)
	fmt.Printf("%s\n", output.Msg)
	assert.Equal(t, 0, output.Code)

	downinput := io.StorageDownloadInput{
		DataHash: output.DataHash,
	}
	outputd := DowloadDataCore(&downinput, TestOntId)
	assert.Equal(t, 0, outputd.Code)
	assert.Equal(t, upinput.DataSource, outputd.DataSource)
}
