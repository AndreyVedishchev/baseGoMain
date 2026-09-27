package db

import (
	"fmt"
	"os"

	"base-go/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storage struct {
	dbConnection *sqlx.DB
	log          *zap.Logger
}

func NewStorage(log *zap.Logger) (*Storage, error) {
	str := &Storage{}
	str.log = log
	_, err := str.connection()
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (as *Storage) connection() (*sqlx.DB, error) {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	dbname := getEnv("POSTGRES_DB", "base_go")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	var err error
	as.dbConnection, err = sqlx.Connect("postgres", url)
	if err != nil {
		as.log.Error("ошибка подключения к базе данных", zap.Error(err))
		return nil, err
	}
	as.log.Info("успешное подключение к базе данных")
	return as.dbConnection, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (as *Storage) Close() {
	as.dbConnection.Close()
	as.log.Info("соединение с базой данных закрыто")
}

/////////////////////////////////////////////////////////////////////

// Save сохраняет модель в хранилище
func (as *Storage) Save(r *models.Resume) {
	as.log.Debug("вызов функции Save", zap.String("uuid", r.UUID))
	_, err := as.dbConnection.NamedExec("INSERT INTO resumes (uuid) VALUES (:uuid)", r)
	if err != nil {
		as.log.Error("ошибка при сохранении новой записи в БД", zap.String("uuid", r.UUID), zap.Error(err))
	}
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *Storage) Delete(uuid string) {
	as.log.Debug("вызов функции Delete", zap.String("uuid", uuid))
	_, err := as.dbConnection.Exec("delete from resumes where uuid = $1", uuid)
	if err != nil {
		as.log.Error("ошибка при удалении записи по uuid из БД", zap.String("uuid", uuid), zap.Error(err))
	}
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *Storage) Get(uuid string) *models.Resume {
	as.log.Debug("вызов функции Get", zap.String("uuid", uuid))
	var resume models.Resume
	err := as.dbConnection.Get(&resume, "select uuid from resumes where uuid=$1", uuid)
	if err != nil {
		as.log.Error("ошибка поиска записи по uuid из БД", zap.String("uuid", uuid), zap.Error(err))
		return nil
	}
	return &resume
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *Storage) Size() int {
	as.log.Debug("вызов функции Size")
	var count int
	err := as.dbConnection.Get(&count, "select count(uuid) cnt from resumes")
	if err != nil {
		as.log.Error("ошибка при подсчете кол-ва записей в БД", zap.Error(err))
		return 0
	}
	return count
}

// GetAll возвращает набор ненулевых резюме
func (as *Storage) GetAll() []*models.Resume {
	as.log.Debug("вызов функции GetAll")
	var resumes []*models.Resume
	err := as.dbConnection.Select(&resumes, "select uuid from resumes")
	if err != nil {
		as.log.Error("ошибка при чтении всех записей из БД", zap.Error(err))
		return nil
	}
	return resumes
}

// Clear удаляет все элементы из хранилища
func (as *Storage) Clear() {
	as.log.Debug("вызов функции Clear")
	_, err := as.dbConnection.Exec("delete from resumes")
	if err != nil {
		as.log.Error("ошибка при очистке БД", zap.Error(err))
	}
}
