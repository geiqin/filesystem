package filesystem

import (
	"mime/multipart"
)


type WhoStorage interface {
	Upload(fileInfo *FileInfo, fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error)
}
