package postgresql

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New("postgres://postgres:12@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStore_Posts(t *testing.T) {
	posts := []Post{
		{
			Title: "Test Post",
			Link:  strconv.Itoa(rand.Intn(1_000_000_000)),
		},
	}
	db, err := New("postgres://postgres:12@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddPost(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.Posts(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}
