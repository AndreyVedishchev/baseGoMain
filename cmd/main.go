package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"base-go/internal/models"
	"base-go/internal/storage/array"
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
	reader := bufio.NewReader(os.Stdin)
	var arrayStorage storage = array.NewArrayStorage()

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

		switch parts[0] {
		case "list":
			printAll(arrayStorage)
		case "size":
			fmt.Println(arrayStorage.Size())
		case "save":
			arrayStorage.Save(&models.Resume{UUID: uuid})
			printAll(arrayStorage)
		case "delete":
			arrayStorage.Delete(uuid)
			printAll(arrayStorage)
		case "get":
			fmt.Println(arrayStorage.Get(uuid))
		case "clear":
			arrayStorage.Clear()
			printAll(arrayStorage)
		case "exit":
			return
		default:
			fmt.Println("Неверная команда.")
		}
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
