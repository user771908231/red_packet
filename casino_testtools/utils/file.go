package utils

import (
	"mime/multipart"
	"os"
	"io"
)

//复制文件
func SaveFileTo(src *multipart.FileHeader, file_path string, file_name string) error {
	f_src, err := src.Open()
	defer f_src.Close()
	if err != nil {
		return err
	}
	dst, err:= os.OpenFile(file_path + file_name, os.O_WRONLY|os.O_CREATE, 0644)
	defer dst.Close()
	if err != nil {
		return err
	}
	_,err = io.Copy(dst, f_src)
	return nil
}
