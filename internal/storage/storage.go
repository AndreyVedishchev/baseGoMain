package storage

import "base-go/internal/models"

// Storage описывает интерфейс хранилища
type Storage interface {
	Save(r *models.Resume)
	Delete(uuid string)
	Get(uuid string) *models.Resume
	Size() int
	GetAll() []*models.Resume
	Clear()
}
