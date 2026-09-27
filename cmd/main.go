package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"base-go/internal/logger"
	"base-go/internal/metrics"
	"base-go/internal/models"
	"base-go/internal/storage/db"
)

// storage описывает интерфейс хранилища
type storage interface {
	Save(r *models.Resume)
	Delete(uuid string)
	Get(uuid string) *models.Resume
	Size() int
	GetAll() []*models.Resume
	Clear()
}

func main() {
	log, err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if err := godotenv.Load(); err != nil {
		log.Warn("файл .env не найден")
	}

	m := metrics.NewMetrics()

	go startMetricsServer(log)

	reader := bufio.NewReader(os.Stdin)
	pgStorage, err := db.NewStorage(log)
	if err != nil {
		log.Error("не удалось подключиться к БД", zap.Error(err))
	}
	defer pgStorage.Close()

	for {
		fmt.Print("Введите одну из команд - (list | size | save uuid | delete uuid | get uuid | clear | exit): ")
		params, _ := reader.ReadString('\n')
		params = strings.TrimSpace(strings.ToLower(params))
		parts := strings.Split(params, " ")

		if len(parts) < 1 || len(parts) > 2 {
			fmt.Println("Неверная команда.")
			continue
		}

		var uuid string
		if len(parts) == 2 {
			uuid = parts[1]
		}

		command := parts[0]
		start := time.Now()
		status := "ok"

		switch command {
		case "list":
			printAll(pgStorage)
		case "size":
			fmt.Println(pgStorage.Size())
		case "save":
			pgStorage.Save(&models.Resume{UUID: uuid})
			m.StorageSize.Set(float64(pgStorage.Size()))
			printAll(pgStorage)
		case "delete":
			pgStorage.Delete(uuid)
			m.StorageSize.Set(float64(pgStorage.Size()))
			printAll(pgStorage)
		case "get":
			fmt.Println(pgStorage.Get(uuid))
		case "clear":
			pgStorage.Clear()
			m.StorageSize.Set(float64(pgStorage.Size()))
			printAll(pgStorage)
		case "exit":
			return
		default:
			status = "error"
			fmt.Println("Неверная команда.")
		}

		m.CommandsTotal.WithLabelValues(command, status).Inc()
		m.CommandDuration.WithLabelValues(command).Observe(time.Since(start).Seconds())
	}
}

// метрики http://localhost:2112/metrics
// дашборд http://localhost:9090/
func startMetricsServer(log *zap.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":2112", mux); err != nil {
		log.Error("ошибка сервера метрик", zap.Error(err))
	}
}

func printAll(arrayStorage storage) {
	all := arrayStorage.GetAll()
	fmt.Println("----------------------------")
	if len(all) == 0 {
		fmt.Println("Empty")
	} else {
		for _, r := range all {
			fmt.Println(r)
		}
	}
	fmt.Println("----------------------------")
}
