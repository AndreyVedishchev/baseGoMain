package db

import (
	"base-go/internal/logger"
	"base-go/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
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
		logger.Log.Error("ошибка подключения к базе данных", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("успешное подключение к базе данных")
	return dbConnection, nil
}

func CloseConnection() {
	dbConnection.Close()
	logger.Log.Info("соединение с базой данных закрыто")
}

/////////////////////////////////////////////////////////////////////

// Save сохраняет модель в хранилище
func (as *ArrayStorage) Save(r *models.Resume) {
	logger.Log.Debug("вызов функции Save", zap.String("uuid", r.UUID))
	_, err := dbConnection.NamedExec("INSERT INTO resumes (uuid) VALUES (:uuid)", r)
	if err != nil {
		logger.Log.Error("ошибка при сохранении новой записи в БД", zap.String("uuid", r.UUID), zap.Error(err))
	}
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *ArrayStorage) Delete(uuid string) {
	logger.Log.Debug("вызов функции Delete", zap.String("uuid", uuid))
	_, err := dbConnection.Exec("delete from resumes where uuid = $1", uuid)
	if err != nil {
		logger.Log.Error("ошибка при удалении записи по uuid из БД", zap.String("uuid", uuid), zap.Error(err))
	}
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *ArrayStorage) Get(uuid string) *models.Resume {
	logger.Log.Debug("вызов функции Get", zap.String("uuid", uuid))
	var resume models.Resume
	err := dbConnection.Get(&resume, "select * from resumes where uuid=$1", uuid)
	if err != nil {
		logger.Log.Error("ошибка поиска записи по uuid из БД", zap.String("uuid", uuid), zap.Error(err))
		return nil
	}
	return &resume
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *ArrayStorage) Size() int {
	logger.Log.Debug("вызов функции Size")
	var count int
	err := dbConnection.Get(&count, "select count(*) cnt from resumes")
	if err != nil {
		logger.Log.Error("ошибка при подсчете кол-ва записей в БД", zap.Error(err))
		return 0
	}
	return count
}

// GetAll возвращает набор ненулевых резюме
func (as *ArrayStorage) GetAll() []*models.Resume {
	logger.Log.Debug("вызов функции GetAll")
	var resumes []*models.Resume
	err := dbConnection.Select(&resumes, "select * from resumes")
	if err != nil {
		logger.Log.Error("ошибка при чтении всех записей из БД", zap.Error(err))
		return nil
	}
	return resumes
}

// Clear удаляет все элементы из хранилища
func (as *ArrayStorage) Clear() {
	logger.Log.Debug("вызов функции Clear")
	_, err := dbConnection.Exec("delete from resumes")
	if err != nil {
		logger.Log.Error("ошибка при очистке БД", zap.Error(err))
	}
}
