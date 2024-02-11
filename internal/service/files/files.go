package service

import (
	"file-chunker/internal/models"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func downloadFile(url string, destination string) error {
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Downloaded %s to %s", url, destination)

	return nil
}

func (s *FilesService) UploadFile(handler *multipart.FileHeader) (string, <-chan error) {
	var wg sync.WaitGroup
	wg.Add(1)

	chunks := make(chan []byte, 1)
	errors := make(chan error)

	var streamError error

	go func() {
		_, streamError = s.chunking.ChunkStreaming(&wg, handler, chunks)
	}()

	file := models.File{
		UID:  uuid.New().String(),
		Name: handler.Filename,
		Size: handler.Size,
	}

	err := s.filesStorage.SaveFile(&file)
	if err != nil {
		errors <- err
		close(errors)
		return "", errors
	}

	chunkCounter := 0
	go func() {
		for chunk := range chunks {
			wg.Add(1)
			log.Println("Received chunk from streamer:", chunkCounter)

			chunkName := fmt.Sprintf("chunk_%d", chunkCounter)
			msg, uploadErr := s.discord.SendChunk(chunkName, chunk)
			if uploadErr == nil {
				chunkModel := models.Chunks{
					FileUID:    file.UID,
					DiscordURL: msg["attachments"].([]interface{})[0].(map[string]interface{})["url"].(string),
					Name:       chunkName,
					Order:      chunkCounter,
				}
				_ = s.filesStorage.SaveChunk(&chunkModel)
			} else {
				errors <- uploadErr
			}
			chunkCounter++
			wg.Done()
		}
		wg.Wait()
		if streamError != nil {
			errors <- streamError
		}
		close(errors)
	}()

	return file.UID, errors
}

func (s *FilesService) DownloadFile(uid string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	downloadsPath := filepath.Join(homeDir, "Downloads")
	destination := fmt.Sprintf("%s/%s", downloadsPath, uid)
	err = os.Mkdir(destination, 0755)
	if err != nil {
		return err
	}
	defer os.RemoveAll(destination)

	chunks, err := s.filesStorage.GetChunks(uid)
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		err := downloadFile(chunk.DiscordURL, fmt.Sprintf("%s/%s", destination, chunk.Name))
		if err != nil {
			return err
		}
	}

	file, err := s.filesStorage.GetFile(uid)
	if err != nil {
		return err
	}

	err = s.chunking.RestoreFile(destination, file.Name, chunks)
	if err != nil {
		return err
	}

	return beeep.Alert("File downloaded", fmt.Sprintf("File %s downloaded", file.Name), "")
}

func (s *FilesService) Files(limit, offset int) ([]models.File, error) {
	return s.filesStorage.Files(limit, offset)
}

func (s *FilesService) DeleteFile(uid string) error {
	return s.filesStorage.DeleteFile(uid)
}

func (s *FilesService) SearchFile(name string) ([]models.File, error) {
	return s.filesStorage.SearchFile(name)
}
