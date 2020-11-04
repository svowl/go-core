package storage

// ReaderWriter это интерфейс чтения данных из внешнего источника
type ReaderWriter interface {
	ReadAll() ([]byte, error)
	Write([]byte) (int, error)
}
