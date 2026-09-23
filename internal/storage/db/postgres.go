package db

import (
	"base-go/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbConnection *sqlx.DB

// ArrayStorage массив резюме
type ArrayStorage struct {
	resumes []*models.Resume
	cnt     int
}

// NewArrayStorage возвращает новый экземпляр хранилища
func NewArrayStorage() *ArrayStorage {
	return &ArrayStorage{resumes: []*models.Resume{}}
}

func NewConnection() (*sqlx.DB, error) {
	url := "postgres://postgres:postgres@localhost:5432/base_go?sslmode=disable"
	var err error
	dbConnection, err = sqlx.Connect("postgres", url)
	if err != nil {
		fmt.Println("Error connecting to database")
		return nil, err
	}
	fmt.Println("Successfully connected to database")
	return dbConnection, nil
}

func CloseConnection() {
	dbConnection.Close()
	fmt.Println("Successfully closed connection to database")
}

/////////////////////////////////////////////////////////////////////

// Save сохраняет модель в хранилище
func (as *ArrayStorage) Save(r *models.Resume) {
	fmt.Println("Вызов функции Save")
	_, err := dbConnection.NamedExec("INSERT INTO resumes (uuid) VALUES (:uuid)", r)
	if err != nil {
		fmt.Println("Ошибка при сохранении новой записи в БД", err)
	}
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *ArrayStorage) Delete(uuid string) {
	fmt.Println("Вызов функции Delete")
	_, err := dbConnection.Exec("delete from resumes where uuid = $1", uuid)
	if err != nil {
		fmt.Println("Ошибка при удалении записи по uuid из БД", err)
	}
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *ArrayStorage) Get(uuid string) *models.Resume {
	fmt.Println("Вызов функции Get")
	var resume models.Resume
	err := dbConnection.Get(&resume, "select * from resumes where uuid=$1", uuid)
	if err != nil {
		fmt.Println("Ошибка поиска записи по uuid из БД", err)
		return nil
	}
	return &resume
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *ArrayStorage) Size() int {
	fmt.Println("Вызов функции Size")
	var count int
	err := dbConnection.Get(&count, "select count(*) cnt from resumes")
	if err != nil {
		fmt.Println("Ошибка при подсчете кол-ва записей в БД", err)
		return 0
	}
	return count
}

// GetAll возвращает набор ненулевых резюме
func (as *ArrayStorage) GetAll() []*models.Resume {
	fmt.Println("Вызов функции GetAll")
	var resumes []*models.Resume
	err := dbConnection.Select(&resumes, "select * from resumes")
	if err != nil {
		fmt.Println("Ошибка при чтении всех записей из БД", err)
		return nil
	}
	return resumes
}

// Clear удаляет все элементы из хранилища
func (as *ArrayStorage) Clear() {
	fmt.Println("Вызов функции Clear")
	_, err := dbConnection.Exec("delete from resumes")
	if err != nil {
		fmt.Println("Ошибка при очистке БД", err)
	}
}
