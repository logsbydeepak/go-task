package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"syscall"
)

var f *os.File

func LoadFile(filepath string) (*os.File, error) {
	var err error
	f, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func CloseFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func NewWriter() *csv.Writer {
	return csv.NewWriter(f)
}

func NewReader() *csv.Reader {
	return csv.NewReader(f)
}
