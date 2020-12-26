package filesystem

import (
	"github.com/geiqin/gotools/helper"
	"github.com/geiqin/xconfig/client"
	"github.com/geiqin/xconfig/model"
	"log"
	"mime/multipart"
)

type IStorage interface {
	Upload(fileInfo *FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error)
}

//上传文件器
type Uploader struct {
	disk string
	mode string
	conf *model.FileSystemInfo
}

var filesystemCfg *model.FilesystemConfig

func init()  {
	filesystemCfg =client.GetFilesystemConfig()
}

func NewUploader(disk string,mode ...string) *Uploader  {
	if disk ==""{
		disk ="qin_store_public"
	}
	m :="cloud"
	if mode !=nil{
		m2 :=mode[0]
		if m !="cloud" && m !="local"{
			log.Println("错误: Uploader 的 mode 参数只能是 cloud 或者 local 值")
			return nil
		}
		m=m2
	}
	return &Uploader{
		disk:    disk,
		mode:	 m,
	}
}

//上传
func (s *Uploader) Upload(fileHeader *multipart.FileHeader, file multipart.File, path string, fileName ...string) (*FileInfo, error) {
	s.getLocalConf(s.disk)
	log.Println("filesystem config :",helper.JsonEncode(filesystemCfg))
	log.Println("disk config :",helper.JsonEncode(s.conf))
	maker :=NewMakeFile(s.conf,fileHeader,file,path)

	//是否重命名（否在自动生成）
	if fileName != nil {
		maker.SetFileName(fileName[0])
	}

	info,_:=maker.Output()

	//选择七牛云上传
	qi := NewQiniuStorage(s.conf)
	ret, err := qi.Upload(info, fileHeader, file)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

//本地储存配置
func (s *Uploader) getLocalConf(name string)*model.FileSystemInfo {
	return nil
}

//云储存配置
func (s *Uploader) getCloudConf(name string) *model.FileSystemInfo {
	if filesystemCfg ==nil{
		log.Println("错误: 未加载 FilesystemConfig 配置信息")
		return nil
	}
	v :=filesystemCfg.Clouds[name]
	if v==nil{
		log.Println("错误: FilesystemConfig 中未配置存盘名称：",name)
	}
	s.conf =v
	return v
}

