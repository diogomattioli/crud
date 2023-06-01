package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/diogomattioli/crud/pkg/data"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var enableDbLogging bool = false

type Dummy struct {
	ID    int    `json:"id_dummy,omitempty" gorm:"primaryKey"`
	Title string `json:"title,omitempty"`
	Valid bool   `json:"valid,omitempty"`
}

func (o *Dummy) GetID() int {
	return o.ID
}

func (o *Dummy) ValidateCreate(token string) error {
	if !o.Valid {
		return data.ValidationErrorNew(1, "Error - Not Valid")
	}
	if token != "" {
		return data.ValidationErrorNew(1, fmt.Sprintf("Token - %+v", token))
	}

	return nil
}

func (o *Dummy) ValidateUpdate(old *Dummy, token string) error {
	if !o.Valid {
		return data.ValidationErrorNew(1, "Error - Not Valid")
	}
	if token != "" {
		return data.ValidationErrorNew(1, fmt.Sprintf("Token - %+v", token))
	}
	return nil
}

func (o *Dummy) ValidateDelete(token string) error {
	if !o.Valid {
		return data.ValidationErrorNew(1, "Error - Not Valid")
	}
	if token != "" {
		return data.ValidationErrorNew(1, fmt.Sprintf("Token - %+v", token))
	}
	return nil
}

type SubDummy struct {
	ID    int    `json:"id_subdummy" gorm:"primaryKey"`
	Title string `json:"title"`
	Valid bool   `json:"valid"`
	Dummy int    `json:"id_dummy"`
}

func (o *SubDummy) GetID() int {
	return o.ID
}

func (o *SubDummy) ValidateCreate(token string) error {
	if !o.Valid {
		return data.ValidationErrorNew(1, "Error - Not Valid")
	}
	return nil
}

func (o *SubDummy) ValidateUpdate(old *SubDummy, token string) error {
	if !o.Valid {
		return data.ValidationErrorNew(1, "Error - Not Valid")
	}
	return nil
}

func (o *SubDummy) ValidateDelete(token string) error {
	if !o.Valid {
		return data.ValidationErrorNew(1, "Error - Not Valid")
	}
	return nil
}

type DummyDefault struct {
	data.Validate[*DummyDefault] `json:"-" gorm:"-"`
	DummyDefaultID               int    `json:"id_dummy_default,omitempty" gorm:"primaryKey"`
	Title                        string `json:"title,omitempty"`
}

func (o *DummyDefault) GetID() int {
	return o.DummyDefaultID
}

func setupDb(quantity int) {

	var newLogger logger.Interface
	if enableDbLogging {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			},
		)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Dummy{})
	db.AutoMigrate(&SubDummy{})
	db.AutoMigrate(&DummyDefault{})

	for i := 1; i <= quantity; i++ {
		db.Create(&Dummy{ID: i, Title: fmt.Sprintf("title%v", quantity-i+1), Valid: true})
		db.Create(&SubDummy{ID: i*2 - 1, Title: fmt.Sprintf("subtitle%v", quantity-i+1), Valid: true, Dummy: i})
		db.Create(&SubDummy{ID: i * 2, Title: fmt.Sprintf("subtitle%v", quantity-i+1), Valid: true, Dummy: i})
		db.Create(&DummyDefault{DummyDefaultID: i, Title: fmt.Sprintf("title%v", quantity-i+1)})
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
	router.HandleFunc("/dummy/", Create[*Dummy]).Methods("POST")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}", Retrieve[Dummy]).Methods("GET")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}", Update[*Dummy]).Methods("PATCH")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}", Delete[*Dummy]).Methods("DELETE")

	router.HandleFunc("/dummy/{id_dummy:[0-9]+}/subdummy/", ListSub[SubDummy, Dummy]).Methods("GET")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}/subdummy/", CreateSub[*SubDummy, Dummy]).Methods("POST")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}/subdummy/{id_subdummy:[0-9]+}", RetrieveSub[SubDummy, Dummy]).Methods("GET")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}/subdummy/{id_subdummy:[0-9]+}", UpdateSub[*SubDummy, Dummy]).Methods("PATCH")
	router.HandleFunc("/dummy/{id_dummy:[0-9]+}/subdummy/{id_subdummy:[0-9]+}", DeleteSub[*SubDummy, Dummy]).Methods("DELETE")

	router.HandleFunc("/misconfigured/{id_wrong}/subdummy/", ListSub[SubDummy, Dummy]).Methods("GET")
	router.HandleFunc("/misconfigured/{id_wrong}/subdummy/", CreateSub[*SubDummy, Dummy]).Methods("POST")
	router.HandleFunc("/misconfigured/{id_wrong}/subdummy/{id_wrong}", RetrieveSub[*SubDummy, Dummy]).Methods("GET")
	router.HandleFunc("/misconfigured/{id_wrong}/subdummy/{id_wrong}", UpdateSub[*SubDummy, Dummy]).Methods("PATCH")
	router.HandleFunc("/misconfigured/{id_wrong}/subdummy/{id_wrong}", DeleteSub[*SubDummy, Dummy]).Methods("DELETE")

	router.HandleFunc("/dummy_default/", List[DummyDefault]).Methods("GET")
	router.HandleFunc("/dummy_default/", Create[*DummyDefault]).Methods("POST")
	router.HandleFunc("/dummy_default/{id_dummy_default:[0-9]+}", Retrieve[DummyDefault]).Methods("GET")
	router.HandleFunc("/dummy_default/{id_dummy_default:[0-9]+}", Update[*DummyDefault]).Methods("PATCH")
	router.HandleFunc("/dummy_default/{id_dummy_default:[0-9]+}", Delete[*DummyDefault]).Methods("DELETE")

	router.ServeHTTP(rec, req)

	return rec
}

func formData(values map[string]string) ([]byte, string, error) {

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	for key, value := range values {

		fw, err := writer.CreateFormField(key)
		if err != nil {
			return []byte{}, "", err
		}

		_, err = io.Copy(fw, strings.NewReader(value))
		if err != nil {
			return []byte{}, "", err
		}
	}

	writer.Close()

	return body.Bytes(), writer.FormDataContentType(), nil
}
