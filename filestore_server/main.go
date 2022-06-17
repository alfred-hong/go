package main

import (
	"filestore_server/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandlerFunc("/file/upload/suc", handler.UpLoadSucHandler)

	http.HandlerFunc("/file/meta", handler.GetFileMetaHandle)
	http.HandlerFunc("hello", handler.UploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server,Err:%S", err.Error())
	}
}
