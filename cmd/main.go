package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	go startMetricsServer()

	reader := bufio.NewReader(os.Stdin)
	var arrayStorage storage = db.NewArrayStorage()

	_, err := db.NewConnection()
	if err != nil {
		logger.Log.Error("не удалось подключиться к БД", zap.Error(err))
	}
	defer db.CloseConnection()

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
			printAll(arrayStorage)
		case "size":
			fmt.Println(arrayStorage.Size())
		case "save":
			arrayStorage.Save(&models.Resume{UUID: uuid})
			metrics.StorageSize.Set(float64(arrayStorage.Size()))
			printAll(arrayStorage)
		case "delete":
			arrayStorage.Delete(uuid)
			metrics.StorageSize.Set(float64(arrayStorage.Size()))
			printAll(arrayStorage)
		case "get":
			fmt.Println(arrayStorage.Get(uuid))
		case "clear":
			arrayStorage.Clear()
			metrics.StorageSize.Set(float64(arrayStorage.Size()))
			printAll(arrayStorage)
		case "exit":
			return
		default:
			status = "error"
			fmt.Println("Неверная команда.")
		}

		metrics.CommandsTotal.WithLabelValues(command, status).Inc()
		metrics.CommandDuration.WithLabelValues(command).Observe(time.Since(start).Seconds())
	}
}

// метрики http://localhost:2112/metrics
// дашборд http://localhost:9090/
func startMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":2112", mux); err != nil {
		logger.Log.Error("ошибка сервера метрик", zap.Error(err))
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
