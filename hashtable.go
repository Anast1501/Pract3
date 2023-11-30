package main

import (
	"errors"
)

//Структура HashMap представляет хэш-таблицу с фиксированным размером массива table
type HashMap struct { //HashMap представляет хеш-таблицу с массивом table размером 512 элементов
	table [512]*Pair
}

//Структура Pair представляет пару ключ-значение, которую храним в таблице
type Pair struct {
	key   string
	value string
}

//Функция calcHash вычисляет хеш ключа и возвращает индекс в массиве table   (прямое хеширование потому что хеш считается напрямую от ключа)
//Если ключ пуст, то она возращает ошибку                                      (вычисление хэша) Хэш - адрес в таблице
func calcHash(key string, size int) (int, error) {
	if len(key) == 0 {
		return 0, errors.New("no value")
	}
	hash := 0
	for i := 0; i < len(key); i++ { //проход по каждому символу в строке key
		hash += int(key[i])
	}
	return hash % size, nil
}

//Метод Insert вставляет пару ключ-значение в хеш-таблицу
//Если ключ уже существует, то он возвращает ошибку
func (hmap *HashMap) Insert(key string, value string) error {
	p := &Pair{key, value}
	hash, err := calcHash(key, len(hmap.table))
	if err != nil {
		return errors.New("unacceptable key")
	}
	//Проверка наличия элемента в ячейке
	if hmap.table[hash] == nil {
		hmap.table[hash] = p
		return nil
	}
	if hmap.table[hash].key == key {
		return errors.New("this key already exists")
	}
	for i := (hash + 1) % len(hmap.table); i != hash; i = (i + 1) % len(hmap.table) {
		if hmap.table[i] == nil {
			hmap.table[i] = p
			return nil
		}
		if hmap.table[i].key == key {
			return errors.New("this key already exists")
		}
	}
	return errors.New("table is full")
}

//Метод Get получает значение по ключу из хеш-таблицы
//Если ключ не найден, то он возвращает ошибку
func (hmap *HashMap) HGet(key string) (string, error) {
	hash, err := calcHash(key, len(hmap.table))
	if err != nil {
		return "", errors.New("unacceptable key")
	}
	//Поиск значения в ячейке с хэшем
	if hmap.table[hash] != nil && hmap.table[hash].key == key {
		return hmap.table[hash].value, nil
	}
	for i := (hash + 1) % len(hmap.table); i != hash; i = (i + 1) % len(hmap.table) {
		if hmap.table[i] != nil && hmap.table[i].key == key {
			return hmap.table[i].value, nil
		}
	}
	return "", errors.New("no such key")
}

//Метод Del удаляет элемент из хеш-таблицы по ключу
//Если ключ не найден, то он возвращает ошибку
func (hmap *HashMap) HDel(key string) error {
	hash, err := calcHash(key, len(hmap.table))
	if err != nil {
		return errors.New("unacceptable key")
	}
	//Проверка наличия элемента в ячейке с хэшом
	if hmap.table[hash] == nil {
		return errors.New("nothing to delete")
	}
	//Поиск элемента и удаление
	if hmap.table[hash].key == key {
		hmap.table[hash] = nil
		return nil
	}
	for i := (hash + 1) % len(hmap.table); i != hash; i = (i + 1) % len(hmap.table) {
		if hmap.table[i] != nil && hmap.table[i].key == key {
			hmap.table[i] = nil
			return nil
		}
	}
	return errors.New("no such key")
}
