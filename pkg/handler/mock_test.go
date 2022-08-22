package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dummy struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Valid bool   `json:"valid"`
}

func (o Dummy) IsValidCreate() bool {
	return o.Valid
}

func (o Dummy) IsValidUpdate(old Dummy) bool {
	return o.Valid
}

func (o Dummy) IsValidDelete() bool {
	return o.Valid
}

type MockAuth struct {
	shouldFail bool
}

func (a MockAuth) Authenticate(user string, pass string) bool {
	return !a.shouldFail
}

func (a MockAuth) Create() string {
	return "123-token"
}

func (a MockAuth) Use(token string) bool {
	return token == "123-token"
}

func setupDb(quantity int) {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Dummy{})

	for i := 1; i <= quantity; i++ {
		db.Create(&Dummy{ID: i, Title: fmt.Sprintf("title%v", quantity-i+1), Valid: true})
	}

	SetDatabase(db)
}

func destroyDb() {
	db, err := db.DB()
	if err == nil {
		db.Close()
	}
}

func serveHTTP(req *http.Request) *httptest.ResponseRecorder {

	rec := httptest.NewRecorder()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/dummy/", List[Dummy]).Methods("GET")
	router.HandleFunc("/dummy/", Create[Dummy]).Methods("POST")
	router.HandleFunc("/dummy/{id}", Retrieve[Dummy]).Methods("GET")
	router.HandleFunc("/dummy/{id}", Update[Dummy]).Methods("PATCH")
	router.HandleFunc("/dummy/{id}", Delete[Dummy]).Methods("DELETE")
	router.ServeHTTP(rec, req)

	return rec
}

func serveHTTPAuth(req *http.Request) *httptest.ResponseRecorder {

	rec := httptest.NewRecorder()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login/", Login).Methods("POST")
	subrouter := router.PathPrefix("/auth").Subrouter()
	subrouter.Use(Auth)
	subrouter.HandleFunc("/dummy/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	router.ServeHTTP(rec, req)

	return rec
}

func formData(values map[string]string) ([]byte, string, error) {

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	for key, value := range values {

		fw, err := writer.CreateFormField(key)
		if err != nil {
			return []byte{}, "", nil
		}

		_, err = io.Copy(fw, strings.NewReader(value))
		if err != nil {
			return []byte{}, "", nil
		}
	}

	writer.Close()

	return body.Bytes(), writer.FormDataContentType(), nil
}
