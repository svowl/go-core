package mem

// Spider type
type Spider struct{}

// New возвращает новый объект Spider
func New() *Spider {
	s := Spider{}
	return &s
}

// Scan выдает фиктивные результаты индексирования
func (*Spider) Scan(string, int) (map[string]string, error) {
	return map[string]string{
		"https://go.dev":                        "go.dev",
		"https://go.dev/":                       "go.dev",
		"https://go.dev/about":                  "About - go.dev",
		"https://go.dev/learn":                  "Learn - go.dev",
		"https://go.dev/solutions":              "Why Go - go.dev",
		"https://go.dev/solutions#case-studies": "Why Go - go.dev",
		"https://go.dev/solutions#use-cases":    "Why Go - go.dev",
		"https://www.google.com":                "Google",
	}, nil
}
