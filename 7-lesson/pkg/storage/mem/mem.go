package mem

// Storage реализует storage.ReaderWriter
type Storage struct {
	Content []byte
	Error   error
}

// ReadAll отдает данные, сохраненные ранее в Content
func (s *Storage) ReadAll() ([]byte, error) {
	if s.Error != nil {
		return s.Content, s.Error
	}
	return s.Content, nil
}

// Write сохраняет b в памяти
func (s *Storage) Write(b []byte) (n int, err error) {
	s.Content = make([]byte, len(b))
	copy(s.Content, b)
	return len(s.Content), nil
}
