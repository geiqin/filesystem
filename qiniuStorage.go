package filesystem

import (
	"context"
	"github.com/geiqin/xconfig/model"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"log"
	"mime/multipart"
)

type QiniuPutRet struct {
	Hash         string `json:"hash"`
	PersistentID string `json:"persistentId"`
	Key          string `json:"key"`
	//Size         int64  `json:"size"`
	//Bucket       string `json:"bucket"`
	//Name         string `json:"name"`
}

type QiniuStorage struct {
	driverConf *model.FileSystemInfo
}

func NewQiniuStorage(cnf *model.FileSystemInfo) *QiniuStorage {
	return &QiniuStorage{driverConf: cnf}
}

//七牛云图片上传
func (q *QiniuStorage) Upload(fileInfo *FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error) {
	var reader io.Reader = file
	var size = fileHeader.Size

	putPolicy := storage.PutPolicy{
		Scope: q.driverConf.Bucket,
	}
	mac := qbox.NewMac(q.driverConf.AccessKey, q.driverConf.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	formUploader := storage.NewFormUploader(&cfg)

	ret := &QiniuPutRet{} //ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": fileHeader.Filename,
		},
	}

	err := formUploader.Put(context.Background(), &ret, upToken, fileInfo.SaveUrl, reader, size, &putExtra)

	if err != nil {
		log.Println("qiniu put :", err)
		return nil, err
	}
	fileInfo.PersistentId = ret.PersistentID
	fileInfo.Hash = ret.Hash
	fileInfo.Size = size

	return fileInfo, nil
}
