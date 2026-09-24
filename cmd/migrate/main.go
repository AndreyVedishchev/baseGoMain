package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("direction", "up", "направление миграции: up | down")
	flag.Parse()

	url := "postgres://postgres:postgres@localhost:5432/base_go?sslmode=disable"

	m, err := migrate.New("file://migrations", url)
	if err != nil {
		log.Fatalf("не удалось инициализировать миграции: %v", err)
	}

	switch *direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("неизвестное направление: %s (используйте up или down)", *direction)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("миграция завершилась с ошибкой: %v", err)
	}

	fmt.Println("миграции применены успешно")
}