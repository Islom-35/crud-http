package rest

import (
	"context"
	"encoding/json"
	"errors"
	"golang-project-template/internal/domain"
	"log"
	"strconv"

	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Books interface {
	Create(ctx context.Context, book domain.Book) error
	GetById(ctx context.Context, id int64) (domain.Book, error)
	UpdateById(ctx context.Context, id int64, inp domain.UpdateBookInput) error
	GetAll(ctx context.Context)([]domain.Book,error)
	Delete(ctx context.Context,id int64)error
}

type Handler struct {
	booksService Books
}

func NewHandler(books Books) *Handler {
	return &Handler{
		booksService: books,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleWare)
	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("/{id:[0-9]+}", h.getBookById).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}",h.updateBookById).Methods(http.MethodPut)
		books.HandleFunc("",h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}",h.deleteBook).Methods(http.MethodDelete)
	}
	return r
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book domain.Book

	if err = json.Unmarshal(reqBytes, &book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.booksService.Create(context.TODO(), book)

	if err != nil {
		log.Println("createBook() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getBookById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("getBookById() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.booksService.GetById(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getBookById() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(book)
	if err != nil {
		log.Println("getBookByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Context-type", "application/json")
	w.Write(response)

}

func (h *Handler) updateBookById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqBytes, err := io.ReadAll(r.Body)
	if err !=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var inp domain.UpdateBookInput
	if err = json.Unmarshal(reqBytes,&inp); err !=nil{
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	err = h.booksService.UpdateById(context.TODO(),id,inp)
	if err !=nil{
		log.Println("error:",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request){
	books,err :=h.booksService.GetAll(context.TODO())
	if err !=nil{
		log.Println("getAllBooks() error:",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response,err :=json.Marshal(books)
	if err !=nil{
		log.Println("getAllBooks() error:",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Context-Type","application/json")
	w.Write(response)
}

func (h *Handler)deleteBook(w http.ResponseWriter,r *http.Request){
	id, err :=getIdFromRequest(r)
	if err !=nil{
		log.Println("deleteBook() error:",err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.booksService.Delete(context.TODO(),id)
	if err!=nil{
		log.Println("deleteBook() error:",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id cant be 0")
	}

	return id, nil
}
