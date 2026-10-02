package db

import (
	"time"

	"base-go/internal/models"
	"base-go/internal/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storage struct {
	dbConnection *sqlx.DB
	log          *zap.Logger
}

func NewStorage(log *zap.Logger) (*Storage, error) {
	initMetrics()

	str := &Storage{}
	str.log = log
	_, err := str.connection()
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (as *Storage) connection() (*sqlx.DB, error) {
	url := utils.BuildPostgresURL()
	maxRetries := utils.GetEnvInt("DB_MAX_RETRIES", 5)
	retryDelay := utils.GetEnvDuration("DB_RETRY_DELAY", 2*time.Second)

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		as.dbConnection, err = sqlx.Connect("postgres", url)
		if err == nil {
			as.log.Info("успешное подключение к базе данных")
			return as.dbConnection, nil
		}
		as.log.Warn("не удалось подключиться к базе данных, повтор попытки",
			zap.Int("попытка", attempt),
			zap.Int("всего попыток", maxRetries),
			zap.Error(err))
		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	as.log.Error("ошибка подключения к базе данных после всех попыток", zap.Error(err))
	return nil, err
}

func (as *Storage) Close() {
	as.dbConnection.Close()
	as.log.Info("соединение с базой данных закрыто")
}

/////////////////////////////////////////////////////////////////////

// Save сохраняет модель в хранилище
func (as *Storage) Save(r *models.Resume) {
	as.log.Debug("вызов функции Save", zap.String("uuid", r.UUID))
	start := time.Now()
	_, err := as.dbConnection.NamedExec("INSERT INTO resumes (uuid) VALUES (:uuid)", r)
	dbOperationDuration.WithLabelValues("save").Observe(time.Since(start).Seconds())
	status := "ok"
	if err != nil {
		status = "error"
		as.log.Error("ошибка при сохранении новой записи в БД", zap.String("uuid", r.UUID), zap.Error(err))
	}
	dbOperationsTotal.WithLabelValues("save", status).Inc()
}

// Delete удаляет элемент из хранилища (при наличии)
func (as *Storage) Delete(uuid string) {
	as.log.Debug("вызов функции Delete", zap.String("uuid", uuid))
	start := time.Now()
	_, err := as.dbConnection.Exec("delete from resumes where uuid = $1", uuid)
	dbOperationDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	status := "ok"
	if err != nil {
		status = "error"
		as.log.Error("ошибка при удалении записи по uuid из БД", zap.String("uuid", uuid), zap.Error(err))
	}
	dbOperationsTotal.WithLabelValues("delete", status).Inc()
}

// Get возвращает пустую модель, либо модель из хранилища (при наличии)
func (as *Storage) Get(uuid string) *models.Resume {
	as.log.Debug("вызов функции Get", zap.String("uuid", uuid))
	start := time.Now()
	var resume models.Resume
	err := as.dbConnection.Get(&resume, "select uuid from resumes where uuid=$1", uuid)
	dbOperationDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
	if err != nil {
		dbOperationsTotal.WithLabelValues("get", "error").Inc()
		as.log.Error("ошибка поиска записи по uuid из БД", zap.String("uuid", uuid), zap.Error(err))
		return nil
	}
	dbOperationsTotal.WithLabelValues("get", "ok").Inc()
	return &resume
}

// Size возвращает количество ненулевых элементов в хранилище
func (as *Storage) Size() int {
	as.log.Debug("вызов функции Size")
	start := time.Now()
	var count int
	err := as.dbConnection.Get(&count, "select count(uuid) cnt from resumes")
	dbOperationDuration.WithLabelValues("size").Observe(time.Since(start).Seconds())
	if err != nil {
		dbOperationsTotal.WithLabelValues("size", "error").Inc()
		as.log.Error("ошибка при подсчете кол-ва записей в БД", zap.Error(err))
		return 0
	}
	dbOperationsTotal.WithLabelValues("size", "ok").Inc()
	return count
}

// GetAll возвращает набор ненулевых резюме
func (as *Storage) GetAll() []*models.Resume {
	as.log.Debug("вызов функции GetAll")
	start := time.Now()
	var resumes []*models.Resume
	err := as.dbConnection.Select(&resumes, "select uuid from resumes")
	dbOperationDuration.WithLabelValues("get_all").Observe(time.Since(start).Seconds())
	if err != nil {
		dbOperationsTotal.WithLabelValues("get_all", "error").Inc()
		as.log.Error("ошибка при чтении всех записей из БД", zap.Error(err))
		return nil
	}
	dbOperationsTotal.WithLabelValues("get_all", "ok").Inc()
	return resumes
}

// Clear удаляет все элементы из хранилища
func (as *Storage) Clear() {
	as.log.Debug("вызов функции Clear")
	start := time.Now()
	_, err := as.dbConnection.Exec("delete from resumes")
	dbOperationDuration.WithLabelValues("clear").Observe(time.Since(start).Seconds())
	status := "ok"
	if err != nil {
		status = "error"
		as.log.Error("ошибка при очистке БД", zap.Error(err))
	}
	dbOperationsTotal.WithLabelValues("clear", status).Inc()
}
