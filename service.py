from flask import Flask, redirect, send_from_directory, request
import random
import string
import socket



#Выполнение TCP-запросов к серверу (отправка запросов к серверу по протоколу TCP)
def TCPrequest(mainRequest, shortKey, fullLink):
    #Указываем адрес сервера и создаём TCP-соединение
    serverAdress = ('localhost', 6379)
    Connection = socket.create_connection(serverAdress)

    #Формирование строки запроса в зависимости от типа запроса (HPUSH или HGET)
    if mainRequest=='HPUSH':
        fullRequest=f"HPUSH {shortKey} {fullLink}"
    if mainRequest=='HGET':
        fullRequest=f"HGET {shortKey}"

    #Отправка закодированной строки запроса по сокету (в виде байтов) необходимо для передачи данных через сетевое соединение   (отправка к бд)
    Connection.send(fullRequest.encode())

    #Получение ответа от базы данных с использованием буфера размером 512 байт
    GetRequest = Connection.recv(512) #буфер для ответа 

    #Закрытие соединения и возврат полученного ответа
    Connection.close()
    return GetRequest

#Создание объекта Flask
app=Flask(__name__)

#Установка адреса и порта сервера Flask
serverAdress = ('localhost', 7485) #Устанавливаем адрес и порт сервера Flask


#Определение маршрута для главной страницы ("/") с поддержкой методов POST и GET
@app.route('/', methods=['POST', 'GET']) #POST (клиент серверу отправляет) GET(сервер клиенту отправляет)
def main(): 
    shortLink = None        #инциализация переменной для короткой ссылки
    if request.method == 'POST':
        fullLink = request.form['user_input'] #извлечение полной ссылки из формы, отправленной пользователем
        shortKey = ShortKeyGeneration()  #генерация короткого ключа с помощью функции

        #Отправка запроса к базе данных для сохранения соответствия между коротким ключом и полной ссылки
        TCPrequest('HPUSH', shortKey, fullLink)

        #Формирование короткой ссылки для отображения пользователю  и возврат короткой ссылки
        shortLink = f"http://localhost:7485/{shortKey}"
        return shortLink


#Определение маршрута с динамическим параметром shortKey
@app.route('/<shortKey>')
def goToOriginal(shortKey):
    #Выполнение запроса к базе данных с использованием функции TCPrequest для получения полной ссылки
    fullLink = TCPrequest('HGET', shortKey, "")
    #Декодирование полученной полной ссылки   (из байтового в текстоый тип)
    fullLink = fullLink.decode()
    #Перенаправление пользователя на оригинальную полную ссылку
    return redirect(fullLink)

#Определение функции для генерации короткого ключа
def ShortKeyGeneration():
    keyElements = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"  #строка символов из которой будет формироваться короткий ключ
    keySize = random.randint(1, 8)  #генерация случайного размера короткого ключа от 1 до 8 символов

    #Генерация короткого ключа путём выбора случайных символов из строки keyElements
    shortKey = ''.join(random.choice(keyElements) for _ in range(keySize))
    return shortKey   #возврат сгенерированного короткого ключа

#Определение маршрута для обслуживания запроса на иконку сайта
@app.route('/favicon.ico')
def favicon():
    #Возврат иконки сайта из директории static с указанным именем и типом контекста
    return send_from_directory('static', 'favicon.ico', mimetype='image/vnd.microsoft.icon')

#Проверка, является ли скрипт основным (а не импортируемым как модуль)
if __name__ =='__main__':
    #Запуск Flask (фреймворк) на локальном сервере с указанным хостом, портом и отключением режима откладки
    app.run(host='localhost', port=7485, debug=False)