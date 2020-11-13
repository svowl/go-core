package mem

// Storage реализует storage.Interface
type Storage struct {
	Content []byte
	Error   error // служит для эмуляции ошибки
}

// New создает переменную типа Storage
func New() *Storage {
	f := Storage{}
	return &f
}

// Get выдает содержимое переменной Storage.Content
func (s *Storage) Get() ([]byte, error) {
	if s.Error != nil {
		return nil, s.Error
	}
	return s.Content, nil
}

// Set сохраняет строку байт в памяти
func (s *Storage) Set(b []byte) (n int, err error) {
	if s.Error != nil {
		return 0, s.Error
	}
	s.Content = make([]byte, len(b))
	copy(s.Content, b)
	return len(s.Content), nil
}
