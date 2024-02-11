package service

import (
	"file-chunker/internal/service/files/chunker"
	"file-chunker/internal/service/files/discord"
	"file-chunker/internal/storage"
)

type FilesService struct {
	filesStorage *storage.Storage
	chunking     *chunker.FileChunking
	discord      *discord.Discord
}

func New(filesStorage *storage.Storage, chunking *chunker.FileChunking, discord *discord.Discord) *FilesService {
	return &FilesService{
		filesStorage: filesStorage,
		chunking:     chunking,
		discord:      discord,
	}
}
