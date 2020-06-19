package utils

import (
	"strconv"
)

func Int2str(n int) string {
	return strconv.Itoa(n) // s == "97"
	//str 2 int64
}

func Str2Int(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func Int642str(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Str2int64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}


func Uint642str(n uint64) string {
	return strconv.FormatUint(n, 10)
}

func Str2UInt64(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}