package pg

import (
	"filmoteka/pkg/films"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"
)

var testdb *Pg

func TestMain(m *testing.M) {
	// Создаем объект БД
	testdb = New("postgres://postgres:12345@localhost:5432/test")

	// Загружаем тестовый SQL-дамп в БД.
	// Дамп соответствует разработанному в задании к 20 уроку.
	files := []string{
		"./sql/1-schema.sql",
		"./sql/2-data.sql",
		"./sql/3-data.sql",
	}

	for _, f := range files {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := testdb.loadSQL(string(file)); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	os.Exit(m.Run())
}

func TestDb_films(t *testing.T) {
	// Добавляем фильм
	addFilms := []films.Film{
		{
			Title:    "Pulp Fiction",
			Year:     1994,
			Profit:   213928762,
			Pgrating: "PG-18",
		},
	}
	err := testdb.AddFilms(addFilms)
	if err != nil {
		t.Fatalf("AddFilm() error: %v", err)
	}
	// Проверяем добавился ли
	allFilms, err := testdb.Films(0)
	if err != nil {
		t.Fatalf("Films(0) error: %v", err)
	}
	var addedFilmID int
	for _, f := range allFilms {
		if f.Title == "Pulp Fiction" {
			addedFilmID = f.ID
			break
		}
	}
	if addedFilmID == 0 {
		t.Fatal("Фильм 'Pulp Fiction' не добавлен")
	}

	// Изменяем название фильма и проверяем есть ли изменения
	updFilm := addFilms[0]
	updFilm.Title = "Pulp Fiction 2"
	testdb.UpdateFilm(addedFilmID, updFilm)

	allFilms, err = testdb.Films(0)
	if err != nil {
		t.Fatalf("Films(0) error: %v", err)
	}
	for _, f := range allFilms {
		if f.ID != addedFilmID {
			continue
		}
		want := "Pulp Fiction 2"
		got := f.Title
		if got != want {
			t.Fatalf("UpdateFilm(): получено %v, ожидается %v", got, want)
		}
		break
	}

	// Удаляем фильм и проверяем удалился ли
	if err := testdb.DeleteFilm(addedFilmID); err != nil {
		t.Fatalf("DeleteFilm() error: %v", err)
	}
	allFilms, err = testdb.Films(0)
	if err != nil {
		t.Fatalf("Films(0) error: %v", err)
	}
	found := false
	for _, f := range allFilms {
		if f.ID == addedFilmID {
			found = true
			break
		}
	}
	if found == true {
		t.Fatal("Фильм 'Pulp Fiction 2' не удален")
	}
}

func TestDb_Films(t *testing.T) {
	tests := []struct {
		name string
		ID   int
		want []string
	}{
		{
			name: "Поиск всех фильмов",
			ID:   0,
			want: []string{
				"Pretty Woman",
				"My Best Friend's Wedding",
				"Notting Hill",
				"Eyes Wide Shut",
				"Mission: Impossible",
				"Interview with the Vampire: The Vampire Chronicles",
				"Once Upon a Time... in Hollywood",
				"Inglourious Basterds",
				"What's Eating Gilbert Grape",
				"Pirates of the Caribbean: The Curse of The Black Pearl",
			},
		},
		{
			name: "Поиск фильмов студии Paramount",
			ID:   5,
			want: []string{
				"Mission: Impossible",
				"What's Eating Gilbert Grape",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilms, err := testdb.Films(tt.ID)
			if err != nil {
				t.Errorf("Films() error: %v", err)
			}
			sort.Strings(tt.want)
			got := make([]string, 0, 10)
			for _, f := range gotFilms {
				got = append(got, f.Title)
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("%v: получено: %v, ожидается %v", tt.name, got, tt.want)
			}
		})
	}
}
