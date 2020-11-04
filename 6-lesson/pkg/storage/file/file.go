package file

import (
	"io/ioutil"
	"os"
)

// Storage реализует storage.ReaderWriter: чтение и запись в файл
type Storage struct {
	FileName string
	offset   int
}

// ReadAll читает все данные
func (s *Storage) ReadAll() ([]byte, error) {
	var b []byte
	b, err := ioutil.ReadFile(s.FileName)
	if err != nil {
		return b, err
	}
	return b, nil
}

// Write сохраняет строку байт b в файл
func (s *Storage) Write(b []byte) (n int, err error) {
	// Записываем b в файл
	f, err := os.OpenFile(s.FileName, os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err = f.Write(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}
