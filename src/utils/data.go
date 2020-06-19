package utils

import (
	"encoding/gob"
	"os"
)

func MergeMap(oldMap map[interface{}]interface{}, newMap map[interface{}]interface{}) map[interface{}]interface{} {
	for k, v := range oldMap {
		newMap[k] = v
	}
	return newMap
}

func WriteMap(path string, data map[interface{}]interface{}) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	e := gob.NewEncoder(file)
	err = e.Encode(data)
	if err != nil {
		return err
	}
	return file.Close()
}

func ReadMap(path string, data *map[interface{}]interface{}) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	d := gob.NewDecoder(file)
	err = d.Decode(data)
	if err != nil {
		return err
	}
	return file.Close()
}

func AppendMap(path string, data map[interface{}]interface{}) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	var oldMap map[interface{}]interface{}
	d := gob.NewDecoder(file)
	err = d.Decode(&oldMap)
	data = MergeMap(oldMap, data)
	e := gob.NewEncoder(file)
	err = e.Encode(data)
	if err != nil {
		return err
	}
	return 	file.Close()
}

func CreateDataDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteDataDir(path string) error {
	return os.RemoveAll(path)
}
