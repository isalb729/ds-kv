package utils

import "hash/fnv"

func BasicHash(key string) uint32 {
	algorithm := fnv.New32a()
	_, _ = algorithm.Write([]byte(key))
	return algorithm.Sum32()
}

