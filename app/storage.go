package main

import (
	"time"
)

type Value struct {
	data string
	exp  int64
}

type Storage struct {
	m map[string]Value
}

func NewStorage() *Storage {
	return &Storage{
		m: make(map[string]Value),
	}
}
func (storage *Storage) get(key string) (string, bool) {
	v, ok := storage.m[key]
	if !ok || isExpired(v.exp) {
		return "", false
	}
	return v.data, ok
}

func (storage *Storage) set(key, value string, exp int64) {
	var laterTime int64
	if exp == -1 {
		laterTime = -1
	} else {
		laterTime = time.Now().Add(time.Duration(exp) * time.Millisecond).UnixMilli()
	}
	storage.m[key] = Value{data: value, exp: laterTime}
	// fmt.Println(storage, time.Now().UnixMilli())
}

func isExpired(exp int64) bool {
	// fmt.Println(time.UnixMilli(exp).UnixMilli(), "\n", time.Now().UnixMilli(), time.Now().After(time.UnixMilli(exp)))
	if exp != -1 && time.Now().After(time.UnixMilli(exp)) {
		return true
	}
	return false
}
