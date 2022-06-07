package file

import (
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type IFile interface {
	GetFormFile(file string) (*File, error)
}

type File struct {
	File     *multipart.File
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	r        *http.Request
}

func (f *File) SetReq(r *http.Request) {
	f.r = r
}

// GetFormFile Get uploaded file
func (f *File) GetFormFile(file string) (*File, error) {
	fileData, handler, err := f.r.FormFile(file)
	if err != nil {
		return nil, err
	}
	f.Header = handler.Header
	f.Filename = handler.Filename
	f.Size = handler.Size
	f.File = &fileData
	return f, nil
}
