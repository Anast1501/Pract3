package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
)
//Объявление глобальной переменной хэш-таблицы для хранения сокращённых и полных ссылок  (имеет тип данных хэш-таблицы)
var urls = HashMap{} //HashMap представляет собой простую хэш-таблицу для хранения сокращённых и полных ссылок 



// handleForm обрабатывает отправку ссылки, перенаправляет на "/shorten" при POST-запросе.
func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return 
	}
}
	//handleShorten обрабатывает укорочение ссылки (URL)
	func handleShorten(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод запроса некорректен", http.StatusMethodNotAllowed)
			return
		}

		//Получение оригинального URL из формы
		originalURL:=r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "Не указан URL", http.StatusBadRequest)
			return 
	}

		//Генерация короткого ключа и вставка пары (короткий ключ, оригинальный URL) в хэш-таблицу
		shortKey:=generateShortKey()
		err:=urls.Insert(shortKey, originalURL) //добавление хэш-таблицы
		if err!=nil {
			http.Error(w, "Невозможно добавить укороченную ссылку в хэш-таблицу", http.StatusInternalServerError)
			return
		}

		// HostIp:=GetMyIP()
		// shortenedURL:=fmt.Sprintf("/short/%s", shortKey)
	}
		
		//handleRedirect обеспечивает перенаправление для укороченных ссылок URL (обработчик перенаправления для укороченных URL)
		func handleRedirect(w http.ResponseWriter, r *http.Request) {
			//Извлечение короткого ключа из URL-пути
			shortKey:=strings.TrimPrefix(r.URL.Path, "/short/")
			if shortKey=="" {
				http.Error(w, "Короткая ссылка пропущена(отсуствует)", http.StatusBadRequest)
				return
			}

			//Получение оригинального URL по короткому ключу и перенаправление
			originalURL, err:=urls.HGet(shortKey)
			if err!=nil {
				http.Error(w, "Короткая ссылка не найдена", http.StatusNotFound)
				return
		}

		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	}

	//generateShortKey генерирует случайный короткий ключ
	func generateShortKey() string {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		keyLength:=rand.Intn(5)+1

		shortKey:=make([]byte, keyLength)
		for i := range shortKey {
			shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

//Получение своего IP-адреса (для того чтобы вставить в короткую ссылку)  ip-адрес:порт/short/(короткий ключ)
//ИЛИ  GetMyIp возвращает IP-адрес локального сервера (для использования в укороченных URL)
func GetMyIP() net.IP {
	conn,_:=net.Dial("udp", "8.8.8.8:80")
	defer conn.Close() //закрытие соединения 
	localAddr:=conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}


//handleUserInput обрабатывает ввод пользователя для укорочения ссылки URL и возвращает укороченный URL
func handleUserInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод запроса некорректен", http.StatusMethodNotAllowed)
		return
	}

	//Получение оригинального URL из формы ввода пользователя
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "Не указан URL", http.StatusBadRequest)
		return
	}

	//Генерация короткого ключа, вставка в хэш-таблицу и формирование укороченного URL для пользователя
	shortKey := generateShortKey()
	err := urls.Insert(shortKey, originalURL)
	if err != nil {
		http.Error(w, "Невозможно добавить укороченную ссылку в хэш-таблицу", http.StatusInternalServerError)
		return
	}

	//Получение моего IP и формирование укороченного URL для пользователя
	HostIP := GetMyIP()
	shortenedURL := fmt.Sprintf("http://%s:%d/short/%s", HostIP, 3333, shortKey)

	//Возвращаем укороченный (сокращённый) URL адрес пользователю 
	fmt.Fprintf(w, "Ваша укороченная ссылка: %s", shortenedURL)
}


func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/short/", handleRedirect)
	http.HandleFunc("/user-input", handleUserInput) //Новый обработчик для ввода пользователя (для пользовательского ввода)                               

	fmt.Println("URL Shortener is running on: 3333")
	http.ListenAndServe("0.0.0.0:3333", nil)
}
