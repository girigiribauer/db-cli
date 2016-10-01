package db

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func openFileReader(filepath string) (*os.File, *bufio.Reader) {
	filepath, _ = homedir.Expand(filepath)

	fp, err := os.Open(filepath)
	if err != nil {
		log.Println("os.Open: ", err)
		return nil, nil
	}

	reader := bufio.NewReader(fp)
	if reader == nil {
		log.Println("bufio.NewReader error")
	}

	return fp, reader
}

func openFileWriter(filepath string) (*os.File, *bufio.Writer) {
	if !strings.HasSuffix(filepath, ".sql") {
		filepath = filepath + ".sql"
	}
	filepath, _ = homedir.Expand(filepath)

	fp, err := os.Create(filepath)
	if err != nil {
		log.Println("os.Create: ", err)
		return nil, nil
	}

	writer := bufio.NewWriter(fp)
	if writer == nil {
		log.Println("bufio.NewWriter error")
	}

	return fp, writer
}
