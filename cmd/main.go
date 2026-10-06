package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"base-go/internal/logger"
	"base-go/internal/models"
	"base-go/internal/storage"
	"base-go/internal/storage/db"
	"base-go/internal/web"
)

func main() {
	log, err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if err := godotenv.Load(); err != nil {
		log.Warn("файл .env не найден")
	}

	reader := bufio.NewReader(os.Stdin)
	pgStorage, err := db.NewStorage(log)
	if err != nil {
		log.Error("не удалось подключиться к БД", zap.Error(err))
	} else {
		go startMetricsServer(log, pgStorage.Registry())
	}
	defer pgStorage.Close()

	webServer := web.NewServer(pgStorage, log)
	go func() {
		if err := webServer.ListenAndServe(":8080"); err != nil {
			log.Error("ошибка веб-сервера", zap.Error(err))
		}
	}()

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
		switch command {
		case "list":
			printAll(pgStorage)
		case "size":
			fmt.Println(pgStorage.Size())
		case "save":
			pgStorage.Save(&models.Resume{UUID: uuid})
			printAll(pgStorage)
		case "delete":
			pgStorage.Delete(uuid)
			printAll(pgStorage)
		case "get":
			fmt.Println(pgStorage.Get(uuid))
		case "clear":
			pgStorage.Clear()
			printAll(pgStorage)
		case "exit":
			return
		default:
			fmt.Println("Неверная команда.")
		}
	}
}

// метрики http://localhost:2112/metrics
// дашборд http://localhost:9090/
func startMetricsServer(log *zap.Logger, reg *prometheus.Registry) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(":2112", mux); err != nil {
		log.Error("ошибка сервера метрик", zap.Error(err))
	}
}

func printAll(arrayStorage storage.Storage) {
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
