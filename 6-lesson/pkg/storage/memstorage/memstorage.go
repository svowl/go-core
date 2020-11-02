package memstorage

// ReaderWriterMem реализует storage.ReaderWriter
type ReaderWriterMem struct {
	Content []byte
	Error   error
}

// ReadAll отдает данные, сохраненные ранее в Content
func (rw *ReaderWriterMem) ReadAll() ([]byte, error) {
	if rw.Error != nil {
		return rw.Content, rw.Error
	}
	return rw.Content, nil
}

// Write сохраняет b в памяти
func (rw *ReaderWriterMem) Write(b []byte) (n int, err error) {
	rw.Content = make([]byte, len(b))
	copy(rw.Content, b)
	return len(rw.Content), nil
}
