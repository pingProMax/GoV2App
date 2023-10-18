package webapi

import (
	"encoding/gob"
	"os"
	"sync"
)

type KeyValue struct {
	lock     *sync.RWMutex
	data     map[string]string
	filePath string
}

var DB *KeyValue

func init() {
	rwm := new(sync.RWMutex)
	DB = &KeyValue{
		filePath: "./my.db",
		data:     make(map[string]string),
		lock:     rwm,
	}
	DB.loadFromFile()
}

func (kv *KeyValue) Get(key string) string {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	value, _ := kv.data[key]
	return value
}

func (kv *KeyValue) Del(key string) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	delete(kv.data, key)
	kv.saveToFile()
}

func (kv *KeyValue) Clear() {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	kv.data = make(map[string]string)
	kv.saveToFile()
}

func (kv *KeyValue) Set(key string, value string) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[key] = value
	kv.saveToFile()
}

func (kv *KeyValue) saveToFile() error {
	file, err := os.Create(kv.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(kv.data); err != nil {
		return err
	}
	return nil
}

func (kv *KeyValue) loadFromFile() error {
	file, err := os.Open(kv.filePath)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&kv.data); err != nil {
		return err
	}
	return nil
}
