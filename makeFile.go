package filesystem

import (
	"fmt"
	"github.com/geiqin/gotools/helper"
	"github.com/geiqin/xconfig/model"
	"mime/multipart"
	"path"
	"strings"
)

//文件信息
type FileInfo struct {
	Title        string `json:"title"`         //文件标题
	Hash         string `json:"hash"`          //文件哈希: 相同文件Hash值相同
	PersistentId string `json:"persistent_id"` //持久化ID
	Size         int64  `json:"size"`          //文件大小
	Duration     int64  `json:"duration"`      //播放时长
	Type         string `json:"type"`          //媒体类型（image/video/voice/document/zip）
	ExtName      string `json:"ext_name"`      //扩展名称
	RawName      string `json:"raw_name"`      //原始文件名称
	FileName     string `json:"file_name"`     //新的文件名称
	Width        int32  `json:"width"`         //图片宽
	Height       int32  `json:"height"`        //图片高
	Path         string `json:"path"`          //相对路径
	SaveUrl      string `json:"save_url"`      //保存完整路径
	Url          string `json:"url"`           //URL
}

type MakeFile struct {
	conf     *model.FileSystemInfo
	header   *multipart.FileHeader
	file     multipart.File
	fileName string
	path     string
}

func NewMakeFile(conf *model.FileSystemInfo, fileHeader *multipart.FileHeader, file multipart.File, path string) *MakeFile {
	obj := &MakeFile{
		conf:     conf,
		header:   fileHeader,
		file:     file,
		path:     path,
		fileName: "",
	}
	return obj
}

func (s *MakeFile) SetFileName(fileName string) {
	s.fileName = fileName
}

func (s *MakeFile) Output() (*FileInfo, error) {
	if s.fileName == "" {
		s.fileName = s.newFileName()
	}

	extName := s.extName()

	info := &FileInfo{
		Title:    s.header.Filename,
		RawName:  s.header.Filename,
		FileName: s.fileName,
		Path:     s.makePath(s.path),
		Type:     extName,
		ExtName:  extName,
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
