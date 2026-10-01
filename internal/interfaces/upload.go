package interfaces

import "mime/multipart"

type UploadProvide interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
}
