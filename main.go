package main

import (
	"fmt"
	"net/http"

	"github.com/yuzhikuan/filestore-server/handler"
)

func main() {
	// 静态资源处理
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.AssetFS())))
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	// 分块上传
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))

	// 用户相关
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SigninHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err: %s", err.Error())
	}
}
