package storage

// Interface это интерфейс получения/сохранения данных в storage
type Interface interface {
	Get() ([]byte, error)
	Set([]byte) (int, error)
}
