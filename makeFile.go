package filesystem

import (
	"fmt"
	"github.com/geiqin/gotools/helper"
	"github.com/geiqin/xconfig/model"
	"mime/multipart"
	"path"
	"strings"
)


type FileInfo struct {
	Title    string `json:"title"`
	Hash     string `json:"hash"`
	Size     int `json:"size"`
	Type     string `json:"type"`
	RawName  string `json:"raw_name"`
	FileName string `json:"file_name"`
	Path     string `json:"path"`
	SaveUrl  string `json:"save_url"`
	Url      string `json:"url"`
}


type MakeFile struct {
	conf   *model.FileSystemInfo
	header *multipart.FileHeader
	file   multipart.File
	fileName string
	path string
}

func NewMakeFile(conf *model.FileSystemInfo,fileHeader *multipart.FileHeader, file multipart.File, path string) *MakeFile {
  obj :=&MakeFile{
	  conf:     conf,
	  header:   fileHeader,
	  file:     file,
	  path: 	path,
	  fileName: "",
  }
  return obj
}

func (s *MakeFile)  SetFileName(fileName string)  {
	s.fileName =fileName
}

func (s *MakeFile) Output() (*FileInfo, error) {
	if s.fileName ==""{
		s.fileName =s.newFileName()
	}

	info := &FileInfo{
		Title:    s.header.Filename,
		RawName:  s.header.Filename,
		FileName: s.fileName,
		Path:     s.makePath(s.path),
		Type:     s.extName(),
	}
	info.SaveUrl = info.Path + "/" + info.FileName
	info.Url = s.makeUrl(info.SaveUrl)

	return info, nil
}


func (s *MakeFile) extName() string {
	ext := path.Ext(s.header.Filename)
	extName, _ := helper.Substr(ext, 1, len(ext))
	if extName != "" {
		extName = strings.ToLower(extName)
	}
	return extName
}

func (s *MakeFile) newFileName() string {
	ext := path.Ext(s.header.Filename)
	newFileName := fmt.Sprintf("%s%s", helper.MD5(helper.UniqueId()), ext)
	return newFileName
}

func (s *MakeFile) makePath(pathStr string) string {
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

func (s *MakeFile) makeUrl(saveUrl string) string {
	url := fmt.Sprintf("%s://%s/%s", s.conf.Transport, s.conf.Domain, saveUrl)
	return strings.ToLower(url)
}

