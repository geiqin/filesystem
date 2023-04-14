package filesystem

import (
	"github.com/geiqin/duration/audio/mp3"
	"github.com/geiqin/duration/audio/wav"
	"github.com/geiqin/duration/video"
	"github.com/geiqin/gotools/helper"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"mime/multipart"
	"time"
)

//媒体处理
type MediaHandle struct {
}

//获取秒部分时间
func (b MediaHandle) GetSeconds(dur time.Duration) int64 {
	t := math.Ceil(dur.Seconds())
	return int64(t)
}

//处理视频和音频时长
func (b MediaHandle) GetVideoAndVoiceDuration(file *multipart.FileHeader) (int64, error) {
	extName := GetFileExtName(file.Filename)
	fileType := GetMediaType(file.Filename)
	f, err := file.Open()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if fileType == "video" {
		mp4len, _ := video.GetMP4Duration(f)
		if mp4len > 0 {
			duration, _ := helper.ToInt64(mp4len)
			return duration, err
		}
	}
	if fileType == "voice" {
		if extName == "mp3" {
			dur, err := mp3.NewDecoder(f).Duration()
			if err != nil {
				log.Println("get mp3 duration", err)
				return 0, err
			}
			return b.GetSeconds(dur), nil
		}
		if extName == "wav" {
			dur, err := wav.NewDecoder(f).Duration()
			if err != nil {
				log.Println("get wav duration", err)
				return 0, err
			}
			return b.GetSeconds(dur), nil
		}
	}
	return 0, nil
}

//处理图片尺寸
func (b MediaHandle) GetImageScale(file *multipart.FileHeader) (width int64, height int64, err error) {
	f, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	if err != nil {
		log.Println("make image scale:", err.Error())
		return 0, 0, err
	} else {
		width, _ = helper.ToInt64(im.Bounds().Dx())
		height, _ = helper.ToInt64(im.Bounds().Dy())
	}
	/*  必须引用这个三个，否在大部分图片解析错误
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	*/
	return width, height, nil
}
