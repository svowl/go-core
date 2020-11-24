package interfaces

import (
	"testing"
)

func Test_older(t *testing.T) {
	p1 := employee{maturity: 25}
	p2 := employee{maturity: 40}
	p3 := customer{age: 75}
	p4 := customer{age: 38}
	p5 := employee{maturity: 20}

	got := older(p1, p2, p3, p4, p5)
	want := 75
	if got != want {
		t.Fatalf("получено %v, ожидается %v", got, want)
	}
}

func Test_mostOlder(t *testing.T) {
	p1 := employee{maturity: 25}
	p2 := employee{maturity: 40}
	p3 := customer{age: 75}
	p4 := customer{age: 38}
	p5 := employee{maturity: 20}

	p := mostOlder(p1, p2, p3, p4, p5)
	var got int
	want := 75
	if p, ok := p.(employee); ok {
		got = p.maturity
	}
	if p, ok := p.(customer); ok {
		got = p.age
	}
	if got != want {
		t.Fatalf("получено %v, ожидается %v", got, want)
	}
}

// Структура реализует интерфейс io.Writer и использутся для тестирования функции echo
type testWriter struct {
	string
}

func (tw *testWriter) Write(b []byte) (int, error) {
	tw.string = tw.string + string(b)
	return len(b), nil
}

func Test_echo(t *testing.T) {
	var tw testWriter
	err := echo(&tw, "This", 12, " is ", employee{10}, "a", true, " ", 0.5, "message")
	if err != nil {
		t.Fatalf("получена ошибка: %v", err)
	}
	got := tw.string
	want := "This is a message"
	if got != want {
		t.Fatalf("получено: %v, ожидается: %v", got, want)
	}
}
