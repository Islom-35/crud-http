package psql

import (
	"context"
	"database/sql"
	"fmt"
	"golang-project-template/internal/domain"
	"strings"
)

type Books struct {
	db *sql.DB
}

func NewBook(db *sql.DB) *Books {
	return &Books{db}
}

func (b *Books) Create(ctx context.Context, book domain.Book) error {
	_, err := b.db.Exec("INSERT INTO books (title, author, publish_date, rating) values ($1,$2,$3,$4)",
		book.Title, book.Author, book.PublishDate, book.Rating)

	return err
}

func (b *Books) GetById(ctx context.Context, id int64) (domain.Book, error) {
	var book domain.Book
	err := b.db.QueryRow("SELECT id,title, author, publish_date, rating from books where id=$1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
	if err == sql.ErrNoRows {
		return book, domain.ErrBookNotFound
	}

	return book, err
}

func (b *Books) UpdateById(ctx context.Context, id int64, inp domain.UpdateBookInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *inp.Title)
		argId++
	}

	if inp.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, *inp.Author)
		argId++
	}

	if inp.PublishDate != nil {
		setValues = append(setValues, fmt.Sprintf("publish_date=$%d", argId))
		argId++
	}

	if inp.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating=$%d", argId))
		argId++
	}
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("Update books Set %s where id=$%d", setQuery, argId)
	args = append(args, id)
	_, err := b.db.Exec(query, args...)
	return err
}

func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := b.db.Query("SELECT id, title, author, publish_date, rating FROM books")
	if err !=nil{
		return nil,err
	}

	books :=make([]domain.Book,0)
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating); err !=nil{
			return nil, err
		}

		books = append(books, book)
	}
	return books,rows.Err()
}

func (b *Books)Delete(ctx context.Context,id int64) error{
	_, err :=b.db.Exec("DELETE FROM books WHERE id=$1",id)
	return err
}
