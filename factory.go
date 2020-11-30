package filesystem

import (
	"fmt"
	"github.com/geiqin/gotools/helper"
	"mime/multipart"
	"path"
	"strings"
)

type Factory struct {
	conf   *DriverConfig
	header *multipart.FileHeader
	file   multipart.File
}

func NewFactory(cnf *DriverConfig) *Factory {
	return &Factory{conf: cnf}
}



func (s *Factory) Upload(fileHeader *multipart.FileHeader, file multipart.File, path string, fileName ...string) (*FileInfo, error) {
	s.header = fileHeader
	s.file = file
	newName := s.newFileName()
	if fileName != nil {
		newName = fileName[0]
	}

	info := &FileInfo{
		Title:    fileHeader.Filename,
		RawName:  fileHeader.Filename,
		FileName: newName,
		Path:     s.makePath(path),
		Type:     s.extName(),
	}
	info.SaveUrl = info.Path + "/" + info.FileName
	info.Url = s.makeUrl(info.SaveUrl)

	//选择七牛云上传
	qi := NewQiniuStorage(s.conf)
	ret, err := qi.Upload(info, fileHeader, file)
	if err != nil {
		return nil, err
	}
	return ret, nil
}


func (s *Factory) extName() string {
	ext := path.Ext(s.header.Filename)
	extName, _ := helper.Substr(ext, 1, len(ext))
	if extName != "" {
		extName = strings.ToLower(extName)
	}
	return extName
}

func (s *Factory) newFileName() string {
	ext := path.Ext(s.header.Filename)
	newFileName := fmt.Sprintf("%s%s", helper.MD5(helper.UniqueId()), ext)
	return newFileName
}

func (s *Factory) makePath(pathStr string) string {
	if pathStr != "" {
		if strings.HasPrefix("/", pathStr) {
			pathStr = strings.TrimPrefix(pathStr, "/")
		}
		if strings.HasSuffix("/", pathStr) {
			pathStr = strings.TrimSuffix(pathStr, "/")
		}
	}
	return strings.ToLower(pathStr)
}

func (s *Factory) makeUrl(saveUrl string) string {
	url := fmt.Sprintf("%s://%s/%s", s.conf.Transport, s.conf.Domain, saveUrl)
	return strings.ToLower(url)
}
