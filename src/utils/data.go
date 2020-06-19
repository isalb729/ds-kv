package utils

import "os"

func WriteStruct(path) {

}

func ReadStruct() {

}

func WriteToMap() {

}

func ReadFromMap() {

}

func CreateDataDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteDataDir(path string) error {
	return os.RemoveAll(path)
}