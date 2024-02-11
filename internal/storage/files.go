package storage

import "file-chunker/internal/models"

type Files interface {
	SaveFile(file *models.File) error
	SaveChunk(chunk *models.Chunks) error
	GetChunks(fileUID string) ([]models.Chunks, error)
	GetFile(uid string) (*models.File, error)
	Files(limit, offset int) ([]models.File, error)
	DeleteFile(uid string) error
	SearchFile(name string) ([]models.File, error)
}

func (s Storage) SaveChunk(chunk *models.Chunks) error {
	return s.DB.Save(&chunk).Error
}

func (s Storage) SaveFile(file *models.File) error {
	return s.DB.Save(&file).Error
}

func (s Storage) GetChunks(fileUID string) ([]models.Chunks, error) {
	var chunks []models.Chunks
	err := s.DB.Where("file_uid = ?", fileUID).Find(&chunks).Error
	if err != nil {
		return nil, err
	}
	return chunks, err
}

func (s Storage) GetFile(uid string) (*models.File, error) {
	var file models.File
	err := s.DB.Where("uid = ?", uid).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s Storage) Files(limit, offset int) ([]models.File, error) {
	var files []models.File
	err := s.DB.Order("id desc").Offset(offset).Limit(limit).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s Storage) DeleteFile(uid string) error {
	err := s.DB.Where("uid = ?", uid).Delete(&models.File{}).Error
	if err != nil {
		return err
	}
	return s.DB.Where("file_uid = ?", uid).Delete(&models.Chunks{}).Error
}

func (s Storage) SearchFile(name string) ([]models.File, error) {
	var files []models.File
	err := s.DB.Order("id desc").Where("name LIKE ?", "%"+name+"%").Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
