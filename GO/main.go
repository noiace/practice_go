package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	app = kingpin.New("dictchecker", "Проверка словаря на дубликаты.")

	// Флаги
	filePath         = app.Flag("path", "Путь к файлу словаря (обязательно)").Short('p').Required().String()
	deleteDuplicates = app.Flag("delete-duplicates", "Удалить дубликаты без подтверждения").Short('d').Default("false").Bool()
	showDuplicates   = app.Flag("show-duplicates", "Показать дубликаты").Short('s').Default("false").Bool()
	mode             = app.Flag("mode", "Режим работы (create/rewrite)").Short('m').Default("rewrite").Enum("create", "rewrite")
	verbose          = app.Flag("verbose", "Вывести подробную информацию").Short('v').Default("false").Bool()
)

func main() {
	var is_valid bool = true
	// Парсим аргументы
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Проверка расширения файла
	if filepath.Ext(*filePath) != ".txt" {
		is_valid := false
		fmt.Println(is_valid)
		if *verbose {
			fmt.Printf("Ошибка: файл должен иметь расширение .txt")
		}
		return
	}

	// Чтение файла
	data, err := os.ReadFile(*filePath)
	if err != nil {
		is_valid := false
		fmt.Println(is_valid)
		if *verbose {
			fmt.Printf("Ошибка чтения файла: %v\n", err)
		}
		return
	}

	if len(data) == 0 {
		is_valid := false
		fmt.Println(is_valid)
		if *verbose {
			fmt.Println("Файл пуст.")
		}
		return
	}

	// Поиск дубликатов
	words := strings.Fields(string(data))
	dict := make(map[string]int)
	for _, word := range words {
		dict[word]++
	}

	// Сбор дубликатов
	var duplicates []string
	for word, count := range dict {
		if count > 1 {
			duplicates = append(duplicates, word)
		}
	}

	if !*showDuplicates && len(duplicates) > 0 {
		is_valid = false
		fmt.Println(is_valid)
		if *verbose {
			fmt.Println("Найдены слова дубликаты")
		}
	}

	// Показ дубликатов (если флаг -s)
	if *showDuplicates && len(duplicates) > 0 {
		fmt.Println("Найденные слова дубликаты:", duplicates)
	}

	// Удаление дубликатов (если флаг -d)
	if *deleteDuplicates && len(duplicates) > 0 {
		uniqueWords := make([]string, 0, len(dict))
		for word := range dict {
			uniqueWords = append(uniqueWords, word)
		}

		outputPath := *filePath
		if *mode == "create" {
			outputPath = strings.TrimSuffix(*filePath, ".txt") + "_unique.txt"
		}

		err = os.WriteFile(outputPath, []byte(strings.Join(uniqueWords, "\n")), 0644)
		if err != nil {
			fmt.Println("Ошибка при записи:", err)
			return
		}

		fmt.Printf("Дубликаты удалены. Результат сохранён в: %s\n", outputPath)
	} else if len(duplicates) == 0 {
		fmt.Println(is_valid)
	}
}
