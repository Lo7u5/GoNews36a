package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

type Post struct {
	Id      int
	Title   string
	Content string
	PubTime int64
	Link    string
}

// New устанавливает новую связь с базой данных
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)

	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

// Posts получение всех записей
func (s *Store) Posts(limit int) ([]Post, error) {
	rows, err := s.db.Query(context.Background(),
		`
			SELECT 
				id,
				title,
				content,
				pub_time,
				link
			FROM posts
			ORDER BY pub_time DESC
			LIMIT $1;
			`, limit)
	if err != nil {
		return nil, err
	}
	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.Id,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// AddPost добавление записи в базу данных
func (s *Store) AddPost(posts []Post) error {
	for _, p := range posts {
		_, err := s.db.Exec(context.Background(),
			`INSERT INTO posts (title, content, pub_time, link) VALUES ($1, $2, $3, $4)`,
			p.Title, p.Content, p.PubTime, p.Link)
		if err != nil {
			return err
		}
	}
	return nil
}
