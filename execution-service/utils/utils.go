package utils

import "encoding/base32"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func BytesToString(bytes []byte) string {
	return base32.StdEncoding.EncodeToString(bytes)
}

func StringToBytes(value string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(value)
}
