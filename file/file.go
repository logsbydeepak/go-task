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

var isNewFile bool

const (
	ID = iota
	DESCRIPTION
	CREATED_AT
	IS_COMPLETE
)

func LoadFile(filepath string) (*os.File, error) {
	_, err := os.Stat(filepath)
	isNewFile = os.IsNotExist(err)

	f, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file")
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

	id, err := strconv.ParseInt(line[ID], 10, 64)
	if err != nil {
		return result, err
	}

	if len(line[DESCRIPTION]) == 0 {
		return result, errors.New("Description can't be empty")
	}

	if len(line[CREATED_AT]) == 0 {
		return result, errors.New("CreatedAt can't be empty")
	}

	isCompleted, err := strconv.ParseBool(line[IS_COMPLETE])
	if err != nil {
		return result, err
	}

	result = Task{
		ID:          id,
		Description: line[DESCRIPTION],
		CreatedAt:   line[CREATED_AT],
		IsComplete:  isCompleted,
	}

	return result, nil
}

func ParseHeader(line []string) error {
	raw := struct {
		id          string
		description string
		createdAt   string
		isComplete  string
	}{
		id:          line[ID],
		description: line[DESCRIPTION],
		createdAt:   line[CREATED_AT],
		isComplete:  line[IS_COMPLETE],
	}

	if raw.id != "ID" {
		return errors.New("Expedited ID")
	}

	if raw.description != "Description" {
		return errors.New("Expedited Description")
	}

	if raw.createdAt != "CreatedAt" {
		return errors.New("Expedited CreatedAt")
	}

	if raw.isComplete != "IsComplete" {
		return errors.New("Expedited IsComplete")
	}

	return nil
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

func Truncate() error {
	err := f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}
