package file

import (
	"io/ioutil"
	"os"
)

// Storage реализует storage.Interface: чтение и запись в файл
type Storage struct {
	FileName string
}

// New создает переменную типа Storage
func New(filename string) *Storage {
	f := Storage{FileName: filename}
	return &f
}

// Get получает данные из файла
func (s *Storage) Get() ([]byte, error) {
	var b []byte
	b, err := ioutil.ReadFile(s.FileName)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Set сохраняет строку байт b в файл
func (s *Storage) Set(b []byte) (n int, err error) {
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
