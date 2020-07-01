package utils

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
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
	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
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

func ReadAllFiles(dir string) ([]string, error) {
	var fileList []string

	dirList := []string{dir}

	for ; len(dirList) != 0; {
		dir = dirList[0]
		dirList = dirList[1:]
		fileInfos, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range fileInfos {
			if file.IsDir() {
				dirList = append(dirList, dir+"/"+file.Name())
			} else {
				fileList = append(fileList, dir+"/"+file.Name())
			}
		}
	}
	return fileList, nil

}

func GetPath(base, key string, storeLevel int) (error, string) {
	if storeLevel < 1 {
		return fmt.Errorf("a store level is supposed to be at least one"), ""
	}
	primes := GetPrimes(storeLevel, 3)
	hash := int(BasicHash(key))
	path := base
	for _, v := range primes {
		path = path + "/" + Int2str(hash%v)
	}
	return nil, path
}

func WriteLocal(data map[string]interface{}, dataDir string, storeLevel int) error {
	for k, v := range data {
		err, path := GetPath(dataDir, k, storeLevel)
		if err != nil {
			return err
		}
		data := map[string]interface{}{}
		err = ReadMap(path, &data)
		if err != nil {
			return err
		}
		err = AppendMap(path, map[string]interface{}{k: v})
		if err != nil {
			return err
		}
	}
	return nil
}
