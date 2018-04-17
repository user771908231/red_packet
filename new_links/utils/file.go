package utils

import (
	"mime/multipart"
	"os"
	"io"
	"errors"
	"bufio"
	"fmt"
	"new_links/model/keysModel"
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

//读取文件
func OpenFiles(string string){
	f,err := os.Open(string)
	if err != nil {
		errors.New("打开文件错误！")
	}
	defer f.Close()
	b := bufio.NewReader(f)
	line, err := b.ReadString(',')
	//arr := []bson.M{}
	for ; err == nil; line, err = b.ReadString(',') {
			K := new(keysModel.Keys)
			K.Keys = line
			K.Upsert()
		//arr = append(arr,str)
	}
	del := os.Remove(string)
	if del != nil {
		fmt.Println(del);
	}

	//if err == io.EOF {
	//	fmt.Print(line)
	//} else {
	//	return errors.New("read occur error!")
	//}
	//return nil
}
