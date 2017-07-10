package utils

import (
	"mime/multipart"
	"os"
	"io"
)

//复制文件
func SaveFileTo(src *multipart.FileHeader, file_path string) error {
	f_src, err := src.Open()
	defer f_src.Close()
	if err != nil {
		return err
	}
	dst, err:= os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0644)
	defer dst.Close()
	if err != nil {
		return err
	}
	_,err = io.Copy(dst, f_src)
	return nil
}

//文件是否存在
func FileIsExist(file_path string) bool {
	_, err := os.Stat(file_path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
