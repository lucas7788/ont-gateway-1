package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/seller_buyer"
	"github.com/zhiqiangxu/ont-gateway/pkg/ddxf/storage/storage"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	key_res_id          = []byte("resourceId")
	key_token_meta_hash = []byte("token_meta_hash")
	key_data_id         = []byte("dataId")
)

func main() {
	db, _ := initDb()
	pwd := []byte("123456")
	wallet, _ := ontology_go_sdk.OpenWallet("./pkg/ddxf/example/wallet.dat")
	seller, _ := wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)
	if false {
		//1.upload
		fileKey, err := upload()
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		//2.SaveDataMeta
		saveDataMetaOutPut, saveDataMetaInput, err := seller_buyer.SaveDataMeta(fileKey)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		//3.SaveTokenMeta
		fmt.Println("DataId: ", saveDataMetaOutPut.DataId)
		db.Put(key_data_id, []byte(saveDataMetaOutPut.DataId), nil)
		saveTokenMetaInput, err := seller_buyer.SaveTokenMeta(saveDataMetaInput.DataMetaHash)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		//4.PublishMeta
		fmt.Println("TokenMetaHash:", saveTokenMetaInput.TokenMetaHash)
		db.Put(key_token_meta_hash, []byte(saveTokenMetaInput.TokenMetaHash), nil)
		resourceId, err := seller_buyer.PublishMeta(seller, saveDataMetaOutPut, saveDataMetaInput, saveTokenMetaInput)
		if err != nil {
			fmt.Println("error: ", err)
		}
		db.Put(key_res_id, []byte(resourceId), nil)
		return
	}

	resourceId, _ := db.Get(key_res_id, nil)
	tokenMetaHash, _ := db.Get(key_token_meta_hash, nil)
	dataId, _ := db.Get(key_data_id, nil)

	buyer, _ := wallet.GetAccountByAddress("AHhXa11suUgVLX1ZDFErqBd3gskKqLfa5N", pwd)
	if true {
		//5.BuyDtoken
		err := seller_buyer.BuyDtoken(buyer, string(resourceId))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		//6.UseToken
		err = seller_buyer.UseToken(buyer, string(resourceId), string(tokenMetaHash), string(dataId))
		if err != nil {
			panic(err)
		}
		return
	}

	if false {
		bs, _ := base64.RawURLEncoding.DecodeString("eyJjb2RlIjowLCJtc2ciOiIiLCJSZXN1bHQiOiJodHRwOi8vMTI3LjAuMC4xOjIwMzM1L2RkeGYvc3RvcmFnZS9kb3dubG9hZC9TdG9yYWdlRmlsZVByZWZpeDA4YmUyZmUzLWY5NTMtNGNkOC1hMTg1LTgyNTEyODhjOGVkOSJ9")
		fmt.Println("bs: ", string(bs))
		return
	}
	if false {
		res, err := upload()
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println("res: ", res)
		return
	}
	if false {
		download2()
		return
	}
}

func initDb() (*leveldb.DB, error) {
	lvlOpts := &opt.Options{
		NoSync: false,
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile("./data", lvlOpts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func download2() {
	url := "http://127.0.0.1:20335/ddxf/storage/download/" + "StorageFilePrefixcb3d9d7b-68dc-42a8-b625-01b0c650e0b9"

	_, _, data, err := forward.Get(url)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("data: ", string(data))
}

func upload() (string, error) {
	url := "http://127.0.0.1:20335/ddxf/storage/upload"
	path := "./pkg/ddxf/example/wallet.dat"
	params := map[string]string{
		"key1": "val1",
	}
	req, err := NewFileUploadRequest(url, path, params)
	if err != nil {
		fmt.Printf("error to new upload file request:%s\n", err.Error())
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error to request to the server:%s\n", err.Error())
		return "", err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	fmt.Println("upload: ", body)
	rr := make(map[string]interface{})
	err = json.Unmarshal(body.Bytes(), &rr)
	if err != nil {
		return "", err
	}
	return rr["fileName"].(string), nil
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
