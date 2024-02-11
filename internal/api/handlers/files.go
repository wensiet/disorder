package handlers

import (
	"file-chunker/internal/service/files"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handlers struct {
	fileService *service.FilesService
}

func New(fs *service.FilesService) *Handlers {
	return &Handlers{
		fileService: fs,
	}
}

func (h *Handlers) UploadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name, errs := h.fileService.UploadFile(file)
		ctx.JSON(http.StatusOK, gin.H{"file": name})
		go func() {
			recv := false
			for range errs {
				beeep.Notify("File uploading", fmt.Sprintf("Failed to load file %s", file.Filename), "")
				recv = true
			}
			if !recv {
				beeep.Notify("File uploading", fmt.Sprintf("File %s has been uploaded", file.Filename), "")
			}
		}()
	}
}

func (h *Handlers) DownloadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileUID := ctx.Param("uid")
		err := h.fileService.DownloadFile(fileUID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"file": fileUID})
	}
}

func (h *Handlers) Files() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
		files, err := h.fileService.Files(limit, offset)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"total": len(files),
			"files": files,
		})
	}
}

func (h *Handlers) DeleteFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("uid")
		err := h.fileService.DeleteFile(uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"file": uid})
	}
}

func (h *Handlers) SearchFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")
		files, err := h.fileService.SearchFile(name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"total": len(files), "files": files})
	}
}
