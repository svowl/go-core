package mem

// Storage реализует storage.IReadWriter
type Storage struct {
	Content []byte
	Error   error // служит для эмуляции ошибки
}

// New создает переменную типа Storage
func New() *Storage {
	f := Storage{}
	return &f
}

// ReadAll выдает содержимое переменной Storage.Content
func (s *Storage) Read() ([]byte, error) {
	if s.Error != nil {
		return nil, s.Error
	}
	return s.Content, nil
}

// Write сохраняет строку байт в памяти
func (s *Storage) Write(b []byte) (n int, err error) {
	if s.Error != nil {
		return 0, s.Error
	}
	s.Content = make([]byte, len(b))
	copy(s.Content, b)
	return len(s.Content), nil
}
