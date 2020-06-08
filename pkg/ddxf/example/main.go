package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/seller_buyer"
	io2 "github.com/zhiqiangxu/ont-gateway/pkg/ddxf/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/storage/storage"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	dataMetaHash  = "e29dd24b167c1e14ea5283aaca199f0da1f9139ace1c5902dfbe20325ff0f61d"
	dataId        = "data_id_12f7d87f-89e6-43fb-b166-b9b103e075bf"
	tokenMetaHash = "e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4"
	qrCodeId      = "publish_id_845f1ea8-634f-4093-852d-8f4a4a7335f2"
	resourceId    = "resource_id_0407bc53-c513-4e24-bb5c-cd64c9279d7f"
)

func main() {

	if false {
		err := upload()
		if err != nil {
			fmt.Println("error: ", err)
		}
		return
	}

	if true {
		download2()
		return
	}

	if false {
		err := seller_buyer.SaveDataMeta()
		if err != nil {
			panic(err)
		}
		return
	}
	if false {
		err := seller_buyer.SaveTokenMeta(dataMetaHash)
		if err != nil {
			panic(err)
		}
		return
	}
	if false {
		err := seller_buyer.PublishMeta1(tokenMetaHash, dataMetaHash)
		if err != nil {
			panic(err)
		}
		return
	}
	if false {
		err := seller_buyer.PublishMeta(qrCodeId)
		if err != nil {
			panic(err)
		}
		return
	}
	pwd := []byte("123456")
	wallet, _ := ontology_go_sdk.OpenWallet("/Users/sss/gopath/src/github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/wallet.dat")
	buyer, _ := wallet.GetAccountByAddress("AHhXa11suUgVLX1ZDFErqBd3gskKqLfa5N", pwd)
	if false {
		err := seller_buyer.BuyDtoken(buyer, resourceId)
		if err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	//012b 64617461   69645f32653231346632372d653539392d346663332d396233662d6336326632666231323436340120e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4
	//012c 64617461 5f69645f63316235663139352d623431342d343535632d393464332d6466303565366563373635300120e2a740fa12bd94f0e242688e29f6d803f7671eb1f81bcfbdc1c3e213878e7dd4
	if true {
		err := seller_buyer.UseToken(buyer, resourceId, tokenMetaHash, dataId)
		if err != nil {
			panic(err)
		}
		return
	}
}

func download2() {
	url := "http://127.0.0.1:20335/ddxf/storage/download"
	input := io2.StorageDownloadInput{
		FileName: "StorageFilePrefix914fb7dc-0ced-47f2-b318-d3a416eeb466",
	}
	bs, _ := json.Marshal(input)
	_, _, data, err := forward.PostJSONRequest(url, bs)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("data: ", string(data))
}

func download() error {
	url := "http://127.0.0.1:20335/ddxf/storage/upload"
	return downloadFile(url, "./download", func(length, downLen int64) {
		fmt.Printf("length: %d, downLen: %d\n", length, downLen)
	})
}
func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	if filesize == info.Size() {
		fmt.Println("安装包已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

func downloadFile(url string, localPath string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	if isFileExist(localPath, fsize) {
		return err
	}
	fmt.Println("fsize", fsize)
	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		fb(fsize, written)
	}
	fmt.Println(err)
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		fmt.Println(err)
	}
	return err
}

func upload() error {
	url := "http://127.0.0.1:20335/ddxf/storage/upload"
	path := "./pkg/ddxf/example/wallet.dat"
	params := map[string]string{
		"key1": "val1",
	}
	req, err := NewFileUploadRequest(url, path, params)
	if err != nil {
		fmt.Printf("error to new upload file request:%s\n", err.Error())
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error to request to the server:%s\n", err.Error())
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println(body)
	return nil
}

func NewFileUploadRequest(url, path string, params map[string]string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(storage.UploadKey, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	// 其他参数列表写入 body
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}
