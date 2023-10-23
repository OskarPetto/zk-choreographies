package files

import (
	"execution-service/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func WritePrivateFile(writeable io.WriterTo, filename string) {
	path := getPrivatePath(filename)
	WriteFile(writeable, path)
}

func ReadPrivateFile(readable io.ReaderFrom, filename string) error {
	path := getPrivatePath(filename)
	return readFile(readable, path)
}

func WritePublicFile(writeable io.WriterTo, filename string) {
	path := getPublicPath(filename)
	WriteFile(writeable, path)
}

func ReadPublicFile(readable io.ReaderFrom, filename string) error {
	path := getPublicPath(filename)
	return readFile(readable, path)
}

func WriteFile(writeable io.WriterTo, path string) {
	file, err := os.Create(path)
	utils.PanicOnError(err)
	defer file.Close()
	bytesWritten, err := writeable.WriteTo(file)
	utils.PanicOnError(err)
	fmt.Printf("Wrote file of size %d in %s\n", bytesWritten, path)
}

func readFile(readable io.ReaderFrom, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	bytesRead, err := readable.ReadFrom(file)
	fmt.Printf("Read %d bytes from %s\n", bytesRead, path)
	return err
}

func getPublicPath(filename string) string {
	return getFilePath() + "public/" + filename
}
func getPrivatePath(filename string) string {
	return getFilePath() + "private/" + filename
}

func getFilePath() string {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "execution-service") {
		wd = filepath.Dir(wd)
	}
	return wd + "/files/"
}
