package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"

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
