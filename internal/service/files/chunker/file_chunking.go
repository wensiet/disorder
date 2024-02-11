package chunker

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"errors"
	"file-chunker/internal/models"
	"file-chunker/internal/service/files/encryptor"
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type FileChunking struct {
	chunkSize int64
	aesKey    []byte
}

func New(chunkSize int64, key []byte) *FileChunking {
	return &FileChunking{
		chunkSize: chunkSize,
		aesKey:    key,
	}
}

func (f *FileChunking) ChunkFile(handler *multipart.FileHeader) (string, error) {
	file, err := handler.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var chunks [][]byte

	for i := 0; i < int(handler.Size); i += int(f.chunkSize) {
		buf := make([]byte, f.chunkSize-aes.BlockSize-md5.Size)
		n, err := file.Read(buf)
		if err != nil {
			return "", err
		}

		var chunk []byte

		encryptedBuf, err := encryptor.EncryptData(buf[:n], f.aesKey)
		if err != nil {
			return "", err
		}
		checksum := md5.Sum(buf[:n])

		chunk = append(chunk, encryptedBuf...)
		chunk = append(chunk, checksum[:]...)

		chunks = append(chunks, chunk)
	}

	bucket := uuid.New().String()

	err = os.Mkdir(bucket, 0755)
	if err != nil {
		return "", err
	}
	for k, v := range chunks {
		outputFilePath := fmt.Sprintf("%s/chunk_%d", bucket, k)
		f, err := os.Create(outputFilePath)
		if err != nil {
			return "", err
		}
		_, err = f.Write(v)
		if err != nil {
			return "", err
		}
		_ = f.Close()
	}
	encryptedMeta, err := encryptor.EncryptData([]byte(handler.Filename), f.aesKey)
	err = os.WriteFile(fmt.Sprintf("%s/metadata", bucket), encryptedMeta, 0644)
	if err != nil {
		return "", err
	}
	return bucket, err
}

func (f *FileChunking) ChunkStreaming(wg *sync.WaitGroup, handler *multipart.FileHeader, output chan []byte) (string, error) {
	defer wg.Done()
	defer close(output)

	file, err := handler.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	for i := 0; i < int(handler.Size); i += int(f.chunkSize) {
		buf := make([]byte, f.chunkSize-aes.BlockSize-md5.Size)
		n, err := file.Read(buf)
		if err != nil {
			return "", err
		}

		var chunk []byte

		encryptedBuf, err := encryptor.EncryptData(buf[:n], f.aesKey)
		if err != nil {
			return "", err
		}
		checksum := md5.Sum(buf[:n])

		chunk = append(chunk, encryptedBuf...)
		chunk = append(chunk, checksum[:]...)
		log.Println("Chunk sent")
		output <- chunk
	}

	return "", nil
}

func (f *FileChunking) RestoreFile(bucket string, filename string, chunks []models.Chunks) error {
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Order < chunks[j].Order
	})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	downloadsPath := filepath.Join(homeDir, "Downloads")

	counter := 0
	resultFilePath := fmt.Sprintf("%s/%s", downloadsPath, filename)

	_, err = os.Stat(resultFilePath)
	for err == nil {
		counter++
		resultFilePath = fmt.Sprintf("%s/%s (%d)", downloadsPath, filename, counter)
		_, err = os.Stat(resultFilePath)
	}

	resultFile, err := os.Create(resultFilePath)
	if err != nil {
		return err
	}
	defer resultFile.Close()

	if err != nil {
		return err
	}
	for _, entry := range chunks {
		chunk, err := os.ReadFile(fmt.Sprintf("%s/%s", bucket, entry.Name))
		if err != nil {
			_ = os.Remove(resultFilePath)
			return err
		}
		decrypted, err := encryptor.DecryptData(chunk[:len(chunk)-md5.Size], f.aesKey)
		if err != nil {
			panic(err)
		}

		realSum := chunk[len(chunk)-md5.Size:]

		calculatedSum := md5.Sum(decrypted)

		if !bytes.Equal(realSum, calculatedSum[:]) {
			return errors.New("checksum mismatch detected")
		}

		_, err = resultFile.Write(decrypted)
		if err != nil {
			_ = os.Remove(resultFilePath)
			return err
		}
	}
	return nil
}
