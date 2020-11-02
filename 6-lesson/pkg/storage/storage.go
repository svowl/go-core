package storage

import (
	"io"
	"os"
)

// ReaderWriter это интерфейс чтения данных из внешнего источника
type ReaderWriter interface {
	ReadAll() ([]byte, error)
	Write([]byte) (int, error)
}

// ReaderWriterFile реализует ReaderWriter: чтение и запись в файл
type ReaderWriterFile struct {
	FileName string
	offset   int
}

// ReadAll читает все данные
func (rw *ReaderWriterFile) ReadAll() ([]byte, error) {
	var b []byte
	rw.offset = 0
	buf := make([]byte, 1024)
	for {
		n, err := rw.read(buf)
		if err == io.EOF {
			b = append(b, buf[:n]...)
			break
		}
		if err != nil {
			return b, err
		}
		b = append(b, buf[:n]...)
	}
	return b, nil
}

// read заполняет буфер b данными из файла
func (rw *ReaderWriterFile) read(b []byte) (n int, err error) {
	f, err := os.Open(rw.FileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err = f.ReadAt(b, int64(rw.offset))
	if err != nil {
		return n, err
	}
	rw.offset += n
	return n, nil
}

// Write сохраняет строку байт b в файл
func (rw *ReaderWriterFile) Write(b []byte) (n int, err error) {
	// Записываем b в файл
	f, err := os.OpenFile(rw.FileName, os.O_CREATE, 0644)
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
