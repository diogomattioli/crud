package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOk(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/", strings.NewReader("{\"id\":0,\"title\":\"title\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var obj Dummy
	db.First(&obj)

	assert.Equal(t, 1, obj.ID)
	assert.Equal(t, "title", obj.Title)
}

func TestCreateEmpty(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("POST", "/dummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec = serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateNotValid(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/", strings.NewReader("{\"id\":0,\"title\":\"title\",\"valid\":false}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateBadJson(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/", strings.NewReader("{\"id\":1,\"title\":\"title_new\":"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateNotJson(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/", strings.NewReader("{\"id\":0,\"title\":\"title\",\"valid\":false}"))
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)

	req.Header.Set("Content-Type", "application/csv")
	rec = serveHTTP(req)

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
}

func TestRetrieveOk(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/7", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var obj Dummy

	err = json.NewDecoder(rec.Body).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 7, obj.ID)
	assert.Equal(t, "title4", obj.Title)
}

func TestRetrieveNotFound(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/7", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRetrieveBadId(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateOk(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader("{\"id\":1,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	var obj Dummy

	db.First(&obj)
	assert.Equal(t, 1, obj.ID)
	assert.Equal(t, "title1", obj.Title)

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	err = json.NewDecoder(rec.Body).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, obj.ID)
	assert.Equal(t, "title_new", obj.Title)
}

func TestUpdateNotJson(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)

	req, err = http.NewRequest("PATCH", "/dummy/1", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/csv")
	rec = serveHTTP(req)

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
}

func TestUpdateEmpty(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("PATCH", "/dummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec = serveHTTP(req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestUpdateBadId(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/a", strings.NewReader("{\"id\":1,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateBadJson(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader("{\"id\":1,\"title\":\"title_new\":"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateInexistent(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader("{\"id\":1,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateInvalid(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader("{\"id\":1,\"title\":\"title_new\",\"valid\":false}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateMismatchingId(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/1", strings.NewReader("{\"id\":2,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteOk(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("DELETE", "/dummy/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteBadId(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("DELETE", "/dummy/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteInexistent(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("DELETE", "/dummy/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteInvalid(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	db.Create(&Dummy{ID: 1, Title: "title", Valid: false})

	req, err := http.NewRequest("DELETE", "/dummy/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
