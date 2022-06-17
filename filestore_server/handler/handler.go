package handler

import (
	"encoding/json"
	"filestore_server/meta"
	"filestore_server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UpLoadAt: time.Now().Format("2006-01-02 15:04:06"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Filed to save date in file,err:%s", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "./file/upload/suc", http.StatusFound)
	}
}

// UpLoadSucHandler:上传已完成
func UpLoadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandle:获取文件元信息
func GetFileMetaHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
