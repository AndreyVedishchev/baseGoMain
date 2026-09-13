package array

import (
	"fmt"

	"base-go/internal/models"
)

// maxSize максимальный размер хранилища
const maxSize = 5

// счетчик занятых элементов в массиве
var cnt int

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
	fmt.Println("Вызов функции Save")
	if cnt < len(as.resumes) {
		as.resumes[cnt] = r
		cnt++
	} else {
		fmt.Println("Элемент не сохранен, кончилось место")
	}
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *ArrayStorage) Delete(uuid string) {
	fmt.Println("Вызов функции Delete")
	for i, v := range as.resumes {
		if v != nil && v.UUID == uuid {
			fmt.Println("найден элемент с uuid:", uuid)
			as.resumes[i] = nil
			cnt--
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
	return cnt
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
	return res
}

// Clear удаляет все элементы из хранилища
func (as *ArrayStorage) Clear() {
	fmt.Println("Вызов функции Clear")
	as.resumes = [maxSize]*models.Resume{}
	cnt = 0
}
