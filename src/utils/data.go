package utils

import (
	"encoding/gob"
	"os"
)

func MergeMap(oldMap *map[string]interface{}, newMap map[string]interface{}) *map[string]interface{} {
	for k, v := range newMap {
		(*oldMap)[k] = v
	}
	return oldMap
}

func WriteMap(path string, data map[string]interface{}) error {
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

func ReadMap(path string, data *map[string]interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	d := gob.NewDecoder(file)
	err = d.Decode(data)
	if err != nil {
		return err
	}
	return file.Close()
}

func AppendMap(path string, data map[string]interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	oldMap := make(map[string]interface{})
	d := gob.NewDecoder(file)
	err = d.Decode(&oldMap)
	if err != nil {
		return err
	}
	// have to write it to cover the original data
	err = file.Close()
	if err != nil {
		return err
	}
	file, _ = os.OpenFile(path, os.O_WRONLY | os.O_CREATE, os.ModePerm)
	pdata := MergeMap(&oldMap, data)
	e := gob.NewEncoder(file)
	err = e.Encode(*pdata)
	if err != nil {
		return err
	}
	return file.Close()
}

func CreateDataDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteDataDir(path string) error {
	return os.RemoveAll(path)
}
