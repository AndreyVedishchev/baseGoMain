package mapp

import (
	"base-go/internal/models"
	"fmt"
)

// maxSize максимальный размер хранилища
const maxSize = 5

// ArrayStorage массив резюме
type ArrayStorage struct {
	resumes map[string]*models.Resume
}

// NewArrayStorage возвращает новый экземпляр хранилища
func NewArrayStorage() *ArrayStorage {
	return &ArrayStorage{resumes: make(map[string]*models.Resume, maxSize)}
}

// Save сохраняет модель в хранилище
func (as *ArrayStorage) Save(r *models.Resume) {
	fmt.Println("Вызов функции Save")
	as.resumes[r.UUID] = r
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *ArrayStorage) Delete(uuid string) {
	fmt.Println("Вызов функции Delete")
	delete(as.resumes, uuid)
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *ArrayStorage) Get(uuid string) *models.Resume {
	fmt.Println("Вызов функции Get")
	res, ok := as.resumes[uuid]
	if ok {
		return res
	}
	return &models.Resume{}
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *ArrayStorage) Size() int {
	fmt.Println("Вызов функции Size")
	fmt.Println("len:", len(as.resumes))
	return len(as.resumes)
}

// GetAll возвращает набор ненулевых резюме
func (as *ArrayStorage) GetAll() []*models.Resume {
	fmt.Println("Вызов функции GetAll")
	resumes := make([]*models.Resume, 0)
	for _, v := range as.resumes {
		if v != nil {
			resumes = append(resumes, v)
		}
	}
	return resumes
}

// Clear удаляет все элементы из хранилища
func (as *ArrayStorage) Clear() {
	fmt.Println("Вызов функции Clear")
	as.resumes = make(map[string]*models.Resume, maxSize)
}
