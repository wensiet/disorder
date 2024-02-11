package api

import (
	"file-chunker/internal/api/handlers"
	service "file-chunker/internal/service/files"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MustRun(service *service.FilesService) {
	r := gin.Default()

	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"*"}
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConf.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConf))

	handling := handlers.New(service)
	r.POST("/upload", handling.UploadFile())
	r.GET("/download/:uid", handling.DownloadFile())
	r.GET("/files", handling.Files())
	r.DELETE("/delete/:uid", handling.DeleteFile())
	r.GET("/search", handling.SearchFile())

	err := r.Run(":8642")
	if err != nil {
		panic(err)
	}
}
