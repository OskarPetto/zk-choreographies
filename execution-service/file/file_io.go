package file

import (
	"fmt"
	"io"
	"os"
	"proof-service/utils"
)

func WriteFile(writeable io.WriterTo, filename string) {
	path := getPath(filename)
	file, err := os.Create(path)
	utils.PanicOnError(err)
	defer file.Close()
	bytesWritten, err := writeable.WriteTo(file)
	utils.PanicOnError(err)
	fmt.Printf("Wrote file of size %d in %s\n", bytesWritten, path)
}

func ReadFile(readable io.ReaderFrom, filename string) error {
	path := getPath(filename)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	bytesRead, err := readable.ReadFrom(file)
	fmt.Printf("Read %d bytes from %s\n", bytesRead, path)
	return err
}

func getPath(filename string) string {
	return "/home/opetto/uni/zk-choreographies/execution-service/files/" + filename
}
