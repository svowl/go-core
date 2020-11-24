package interfaces

import (
	"io"
)

// Задача №1: Обобщение через интерфейс
type employee struct {
	maturity int
}

type customer struct {
	age int
}

func (e employee) years() int {
	return e.maturity
}

func (c customer) years() int {
	return c.age
}

// Обобщающий интерфейс для типов employee и customer
type person interface {
	years() int
}

// Функция принимает на вход список employee и customer, реализующих интерфейс person
// и возвращает возраст самого старшего человека
func older(people ...person) int {
	var maxAge int
	for _, p := range people {
		if p.years() > maxAge {
			maxAge = p.years()
		}
	}
	return maxAge
}

// Задача №2: Обобщение через пустой интерфейс

// Функция принимает список пустых интерфейсов, определяет тип данных и возвращает
// объект с максимальным возрастом
func mostOlder(people ...interface{}) interface{} {
	var maxAge int
	var obj interface{}
	for _, p := range people {
		var age int
		switch p.(type) {
		case employee:
			age = p.(employee).maturity
			break
		case customer:
			age = p.(customer).age
		}
		if age > maxAge {
			maxAge = age
			obj = p
		}
	}
	return obj
}

// Задача №3

// Функция принимает канал вывода w и произвольные аргументы и выводит канал w только строки
func echo(w io.Writer, params ...interface{}) error {
	for _, p := range params {
		if p, ok := p.(string); ok {
			_, err := w.Write([]byte(p))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
