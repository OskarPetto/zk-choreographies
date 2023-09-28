package utils

import "encoding/base64"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func BytesToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func StringToBytes(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}
