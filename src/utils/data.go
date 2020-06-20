package utils

import (
	"encoding/gob"
	"fmt"
	"os"
)

func MergeMap(oldMap *map[string]interface{}, newMap map[string]interface{}) *map[string]interface{} {
	for k, v := range newMap {
		(*oldMap)[k] = v
	}
	return oldMap
}

func ParseDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if i == 0 {
				return "/"
			} else {
				return path[:i]
			}
		}
	}
	return ""
}
func WriteMap(path string, data map[string]interface{}) error {
	dir := ParseDir(path)
	if dir != "" {
		err := CreateDataDir(dir)
		if err != nil {
			return fmt.Errorf("WriteMap create dir: %v", err)
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("WriteMap open file: %v", err)
	}
	defer file.Close()
	e := gob.NewEncoder(file)
	err = e.Encode(data)
	if err != nil {
		return fmt.Errorf("WriteMap encode data: %v", err)
	}
	return nil
}

func ReadMap(path string, data *map[string]interface{}) error {
	dir := ParseDir(path)
	if dir != "" {
		err := CreateDataDir(dir)
		if err != nil {
			return fmt.Errorf("ReadMap create dir: %v", err)
		}
	}
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("ReadMap open file: %v", err)
	}
	defer file.Close()
	d := gob.NewDecoder(file)
	err = d.Decode(data)
	if err != nil {
		if err.Error() == "EOF" {
			return nil
		} else {
			return fmt.Errorf("ReadMap decode data: %v", err)
		}
	}
	return nil
}

func AppendMap(path string, data map[string]interface{}) error {
	dir := ParseDir(path)
	if dir != "" {
		err := CreateDataDir(dir)
		if err != nil {
			return fmt.Errorf("AppendMap create dir: %v", err)
		}
	}
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("AppendMap open file: %v", err)
	}
	oldMap := make(map[string]interface{})
	d := gob.NewDecoder(file)
	err = d.Decode(&oldMap)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("AppendMap decode data: %v", err)
	}
	file.Close()
	file, err = os.OpenFile(path, os.O_WRONLY | os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("AppendMap open file: %v", err)
	}
	defer file.Close()
	pdata := MergeMap(&oldMap, data)
	e := gob.NewEncoder(file)
	err = e.Encode(*pdata)
	if err != nil {
		return fmt.Errorf("AppendMap encode data: %v", err)
	}
	return nil
}

func CreateDataDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteDataDir(path string) error {
	return os.RemoveAll(path)
}

