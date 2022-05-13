package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html界面
		date, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			if _, err := io.WriteString(w, "internal server err"); err != nil {
				panic(err)
			}
		} else {
			if _, err := io.WriteString(w, string(date)); err != nil {
				panic(err)
			}
		}
	} else if r.Method == "PUT" {
		//接受数据流及存储到本地
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get date,err: %s", err.Error())
			return
		}
		defer file.Close()

		newFile, err := os.Create("./tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s", err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Filed to save date in file,err:%s", err.Error())
			return
		}

		http.Redirect(w, r, "./file/upload/suc", http.StatusFound)
	}
}

func UpLoadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
