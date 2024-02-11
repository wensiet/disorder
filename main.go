package main

import (
	"file-chunker/internal/api"
	"file-chunker/internal/config"
	"file-chunker/internal/service/files"
	"file-chunker/internal/service/files/chunker"
	"file-chunker/internal/service/files/discord"
	"file-chunker/internal/storage"
)

var chunkingServ *chunker.FileChunking
var strg *storage.Storage
var filesService *service.FilesService
var ds *discord.Discord

func init() {
	conf := config.GetConfig()

	key, err := config.GetKey()
	if err != nil {
		panic(err)
	}

	strg = storage.New()
	ds = discord.New(conf.Discord.Token, conf.Discord.Channel)

	chunkingServ = chunker.New(conf.Bucket.Size, key)

	filesService = service.New(strg, chunkingServ, ds)

	strg.MustMigrate()
}

func main() {
	api.MustRun(filesService)
}
