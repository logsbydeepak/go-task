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

const filePath = "./task.csv"

var isNewFile bool

func LoadFile() (*os.File, error) {
	_, err := os.Stat(filePath)
	isNewFile = os.IsNotExist(err)

	f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
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

func WriteHeader(writer *csv.Writer) error {
	return writer.Write([]string{"ID", "Description", "CreatedAt", "IsComplete"})
}

func WriteTask(writer *csv.Writer, task Task) error {
	var isComplete string

	if task.IsComplete {
		isComplete = "true"
	} else {
		isComplete = "false"
	}

	return writer.Write([]string{fmt.Sprintf("%v", task.ID), task.Description, task.CreatedAt, isComplete})
}

func IsNewFile() bool {
	return isNewFile
}
