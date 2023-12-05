#Отправка POST-запроса к серверу с целью сокращения введённого URL.
#После отправки запроса, программа выводит сокращённый URL или сообщение об ошибке в зависимости от статуса ответа сервера
import requests

#Определение функции get_short_link, которая отправляет POST-запрос к серверу по указанному URL
def get_short_link(url):
    #Выполнение POST-запроса на сервер по URL 'http://localhost:7485/' с передачей данных в форме  {'user_input': url}
    response = requests.post('http://localhost:7485/', data={'user_input': url})

    #Проверка статус-кода ответа: если код 200, возвращается текст ответа, иначе возвращается сообщение об ошибке
    if response.status_code == 200:
        return response.text
    else:
        return 'Error: ' + response.text

# Пример использования
original_url = str(input("Enter url or command: ")) #получение ввода пользователя в виде URL или команды
short_url = get_short_link(original_url) #вызов функции get_short_link для сокращения введённого URL
print(short_url) #вывод сокращённого URL или сообщения об ошибке