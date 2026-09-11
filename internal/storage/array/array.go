package array

import (
	"fmt"

	"base-go/internal/models"
)

// maxSize максимальный размер хранилища
const maxSize = 10000

// ArrayStorage массив резюме
type ArrayStorage struct {
	resumes [maxSize]*models.Resume
}

// NewArrayStorage возвращает новый экземпляр хранилища
func NewArrayStorage() *ArrayStorage {
	return &ArrayStorage{resumes: [maxSize]*models.Resume{}}
}

// Save сохраняет модель в хранилище
func (as *ArrayStorage) Save(r *models.Resume) {
	//TODO
	fmt.Println("implement me")
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *ArrayStorage) Delete(uuid string) {
	//TODO
	fmt.Println("implement me")
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *ArrayStorage) Get(uuid string) *models.Resume {
	//TODO
	panic("implement me")
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *ArrayStorage) Size() int {
	//TODO
	panic("implement me")
}

// GetAll возвращает набор ненулевых резюме
func (as *ArrayStorage) GetAll() []*models.Resume {
	//TODO
	return []*models.Resume{as.resumes[0]}
}

// Clear удаляет все элементы из хранилища
func (as *ArrayStorage) Clear() {
	//TODO
	fmt.Println("implement me")
}
