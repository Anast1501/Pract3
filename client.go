//Представление собой простого клиента для укорачивания URL-ссылки с помощью сервера (с использованием (отправленных) HTTP-запроса) к серверу
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func main() {
	//Проверка на наличие ровно одного аргумента командной строки
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go <url>") //вывод сообщения о правильном использовании
		os.Exit(1) //завершение работы с кодом выхода 1
	}

	//Получение оригинального URL из аргумента командной строки 
	originalURL := os.Args[1]

	//Вызов функции для получения укороченного URL
	shortenedURL, err := getShortenedURL(originalURL)
	if err != nil { 
		fmt.Printf("Error: %v\n", err) //Вывод сообщения об ошибки и завершение работы с кодом выхода 1
		os.Exit(1)
	}

	//Вывод укороченного URL на экран 
	fmt.Printf("Shortened URL: %s\n", shortenedURL)
}

//Функция getShortenedURL принимает оригинальный URL в качестве параметра и отправляет POST-запрос на сервер
func getShortenedURL(originalURL string) (string, error) {
	//URL сервера для укорачивания  
	url := "http://localhost:3333/user-input" //обработка запроса 

	//Создание тела запроса с передачей оригинального URL
	payload := bytes.NewBufferString(fmt.Sprintf("url=%s", originalURL)) //оформление запроса
	resp, err := http.Post(url, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return "", err //возвращение ошибки в случае неудачного запроса
	}
	defer resp.Body.Close()

	//Проверка успешности запроса по коду состояния 
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned non-OK status: %s", resp.Status) //возвращение ошибки, если код состояния не является 200 ОК	}
	}

	//Чтение укороченного URL из тела ответа 
	var result string
	_, err = fmt.Fscanf(resp.Body, "Ваша укороченная ссылка: %s", &result)
	if err != nil {
		return "", err //возвращение ошибки в случае неудачного чтения из ответа
	}

	return result, nil //возвращение укороченного URL
}
