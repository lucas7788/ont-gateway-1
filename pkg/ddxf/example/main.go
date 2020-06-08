package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ontio/ontology-go-sdk"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
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
	key_res_id          = []byte("resourceId")
	key_token_meta_hash = []byte("token_meta_hash")
	key_data_id         = []byte("dataId")
)

func main() {

	if false {
		bs, _ := base64.RawURLEncoding.DecodeString("eyJjb2RlIjowLCJtc2ciOiIiLCJSZXN1bHQiOiJodHRwOi8vbG9jYWxob3N0L2Jvb2svaGVsbG8ifQ==")
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
	db, _ := initDb()

	pwd := []byte("123456")
	wallet, _ := ontology_go_sdk.OpenWallet("/Users/sss/gopath/src/github.com/zhiqiangxu/ont-gateway/pkg/ddxf/example/wallet.dat")
	seller, _ := wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)
	if false {
		fileKey, err := upload()
		if err != nil {
			fmt.Println("error: ", err)
		}
		saveDataMetaOutPut, saveDataMetaInput, err := seller_buyer.SaveDataMeta(fileKey)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		fmt.Println("DataId: ", saveDataMetaOutPut.DataId)
		db.Put(key_data_id, []byte(saveDataMetaOutPut.DataId), nil)
		saveTokenMetaInput, err := seller_buyer.SaveTokenMeta(saveDataMetaInput.DataMetaHash)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
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
		err := seller_buyer.BuyDtoken(buyer, string(resourceId))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		err = seller_buyer.UseToken(buyer, string(resourceId), string(tokenMetaHash), string(dataId))
		if err != nil {
			panic(err)
		}
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
	fmt.Println(body)
	rr := make(map[string]interface{})
	err = json.Unmarshal(body.Bytes(), rr)
	if err != nil {
		return "", nil
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
