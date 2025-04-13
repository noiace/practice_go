package main

import (
	"fmt"
	"io/ioutil"

	// "os"
	"path/filepath"
	"strings"
)

func main() {
	var filePath string
	is_valid := true

	fmt.Println("Введите полный путь к файлу (например, ./text.txt):")
	fmt.Scanln(&filePath)

	if filepath.Ext(filePath) != ".txt" {
		is_valid = false
		fmt.Println(is_valid)
		fmt.Println("Ошибка: файл должен иметь расширение .txt")
		return
	}

	fmt.Println(filePath)
	data, err := CheckFile(filePath)

	if err != nil {
		is_valid = false
		fmt.Print(is_valid)
		fmt.Println(err)
		return
	}

	if data == "" {
		is_valid = false
		fmt.Println(is_valid)
		fmt.Println("Файл пустой.")
		return
	}

	dict := make(map[string]int)

	for _, value := range strings.Fields(data) {
		dict[value]++
	}

	copyExist := false
	duplicates := []string{}

	for word, count := range dict {
		if count > 1 {
			duplicates = append(duplicates, word)
			copyExist = true
		}
	}

	if copyExist {
		fmt.Println("Найдены дубликаты", duplicates)
		fmt.Println("Желаете ли вы удалить копии y/n?")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) == "y" {

			uniqueWords := make([]string, 0, len(dict))

			for word := range dict {
				uniqueWords = append(uniqueWords, word)
			}

			err := ioutil.WriteFile(filePath, []byte(strings.Join(uniqueWords, "\n")), 0644)

			if err != nil {
				fmt.Println("Ошибка при записи в файл:", err)
				return
			}

			fmt.Println("Дубликаты удалены, файл перезаписан.")
		}

	} else {
		fmt.Println("Словарь валидный. Дубликатов не найдено")
	}

}

func CheckFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл: %v", err)
	}
	return string(data), nil
}
