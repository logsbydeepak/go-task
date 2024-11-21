package file

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
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

type Task struct {
	ID          int64
	Description string
	CreatedAt   string
	IsComplete  bool
}

func ParseLine(line []string) (Task, error) {
	var result Task

	if len(line) != 4 {
		return result, errors.New("Length should be of size 4")
	}

	raw := struct {
		id          string
		description string
		createdAt   string
		isComplete  string
	}{
		id:          line[0],
		description: line[1],
		createdAt:   line[2],
		isComplete:  line[3],
	}

	id, err := strconv.ParseInt(raw.id, 10, 64)
	if err != nil {
		return result, err
	}

	if len(raw.description) == 0 {
		return result, errors.New("Description can't be empty")
	}

	if len(raw.createdAt) == 0 {
		return result, errors.New("Description can't be empty")
	}

	isCompleted, err := strconv.ParseBool(raw.isComplete)
	if err != nil {
		return result, err
	}

	result = Task{
		ID:          id,
		Description: line[1],
		CreatedAt:   line[2],
		IsComplete:  isCompleted,
	}

	return result, nil
}
