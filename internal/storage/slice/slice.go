package slice

import (
	"base-go/internal/models"
	"fmt"
)

// maxSize максимальный размер хранилища
const maxSize = 0

// ArrayStorage массив резюме
type ArrayStorage struct {
	resumes []*models.Resume
}

// NewArrayStorage возвращает новый экземпляр хранилища
func NewArrayStorage() *ArrayStorage {
	return &ArrayStorage{resumes: make([]*models.Resume, 0, maxSize)}
}

// Save сохраняет модель в хранилище
func (as *ArrayStorage) Save(r *models.Resume) {
	fmt.Println("Вызов функции Save")
	as.resumes = append(as.resumes, r)
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *ArrayStorage) Delete(uuid string) {
	fmt.Println("Вызов функции Delete")
	for i, v := range as.resumes {
		if v != nil && v.UUID == uuid {
			as.resumes = append(as.resumes[:i], as.resumes[i+1:]...)
		}
	}
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *ArrayStorage) Get(uuid string) *models.Resume {
	fmt.Println("Вызов функции Get")
	for i, v := range as.resumes {
		if v != nil && v.UUID == uuid {
			fmt.Println("i, v:", i, v)
			return v
		}
	}
	return &models.Resume{}
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *ArrayStorage) Size() int {
	fmt.Println("Вызов функции Size")
	fmt.Println("len:", len(as.resumes), "; cap:", cap(as.resumes))
	return len(as.resumes)
}

// GetAll возвращает набор ненулевых резюме
func (as *ArrayStorage) GetAll() []*models.Resume {
	fmt.Println("Вызов функции GetAll")
	res := []*models.Resume{}
	for _, v := range as.resumes {
		if v != nil {
			res = append(res, v)
		}
	}
	return as.resumes
}

// Clear удаляет все элементы из хранилища
func (as *ArrayStorage) Clear() {
	fmt.Println("Вызов функции Clear")
	as.resumes = make([]*models.Resume, 0, maxSize)
}
