package service

import (
	"context"
	"golang-project-template/internal/domain"
	"time"
)

type BooksRepository interface {
	Create(ctx context.Context, book domain.Book) error
	GetById(ctx context.Context,id int64)(domain.Book,error)
	UpdateById(ctx context.Context,id int64,inp domain.UpdateBookInput)error
	GetAll(ctx context.Context)([]domain.Book,error)
	Delete(ctx context.Context,id int64)error
}

type Books struct{
	repo BooksRepository
}

func NewBooks(repo BooksRepository) *Books{
	return &Books{
		repo: repo,
	}
}

func (b *Books)Create(ctx context.Context, book domain.Book) error{
	if book.PublishDate.IsZero() {
		book.PublishDate = time.Now()
	}
	return b.repo.Create(ctx,book)
}

func (b *Books) GetById(ctx context.Context, id int64)(domain.Book ,error){
	return b.repo.GetById(ctx,id)
}

func (b *Books)UpdateById(ctx context.Context,id int64,inp domain.UpdateBookInput)error{
	return b.repo.UpdateById(ctx,id,inp)
}

func (b *Books) GetAll(ctx context.Context)([]domain.Book, error){
	return b.repo.GetAll(ctx)
}

func(b *Books)Delete(ctx context.Context,id int64)error{
	return b.repo.Delete(ctx,id)
}
