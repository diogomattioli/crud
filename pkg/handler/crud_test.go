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

func TestCreateSubOk(t *testing.T) {

	setupDb(2)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/2/subdummy/", strings.NewReader("{\"id\":0,\"title\":\"title\",\"valid\":true,\"id_dummy\":100}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var slice []SubDummy
	db.Where(SubDummy{Dummy: 2}).Order("id").Find(&slice)

	assert.Equal(t, 3, len(slice))
	assert.Equal(t, 5, slice[2].ID)
	assert.Equal(t, 2, slice[2].Dummy)
	assert.Equal(t, "title", slice[2].Title)
}

func TestCreateSubNotFound(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/dummy/1/subdummy/", strings.NewReader("{\"id\":0,\"title\":\"title\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCreateMisconfigured(t *testing.T) {

	setupDb(2)
	defer destroyDb()

	req, err := http.NewRequest("POST", "/misconfigured/2/subdummy/", strings.NewReader("{\"id\":0,\"title\":\"title\",\"valid\":true,\"id_dummy\":100}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
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

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRetrieveSubOk(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/7/subdummy/14", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var obj SubDummy

	err = json.NewDecoder(rec.Body).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 14, obj.ID)
	assert.Equal(t, 7, obj.Dummy)
	assert.Equal(t, "subtitle4", obj.Title)
}

func TestRetrieveMisconfigured(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/misconfigured/7/subdummy/14", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
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

	assert.Equal(t, http.StatusNotFound, rec.Code)
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

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/6", strings.NewReader("{\"id\":2,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var obj Dummy

	db.First(&obj, 6)
	assert.Equal(t, 6, obj.ID)
	assert.Equal(t, "title_new", obj.Title)

	err = json.NewDecoder(rec.Body).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 6, obj.ID)
	assert.Equal(t, "title_new", obj.Title)
}

func TestUpdateSubOk(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/dummy/5/subdummy/10", strings.NewReader("{\"id\":0,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	var obj SubDummy

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	err = json.NewDecoder(rec.Body).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 10, obj.ID)
	assert.Equal(t, 5, obj.Dummy)
	assert.Equal(t, "title_new", obj.Title)

	db.Where(SubDummy{ID: 10, Dummy: 5}).First(&obj)
	assert.Equal(t, 10, obj.ID)
	assert.Equal(t, 5, obj.Dummy)
	assert.Equal(t, "title_new", obj.Title)
}

func TestUpdateMisconfigured(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("PATCH", "/misconfigured/5/subdummy/10", strings.NewReader("{\"id\":0,\"title\":\"title_new\",\"valid\":true}"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := serveHTTP(req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
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

	assert.Equal(t, http.StatusNotFound, rec.Code)
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

func TestDeleteSubOk(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("DELETE", "/dummy/5/subdummy/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var slice []SubDummy

	db.Where(SubDummy{Dummy: 5}).Find(&slice)
	assert.Equal(t, 1, len(slice))
	assert.Equal(t, 9, slice[0].ID)
	assert.Equal(t, 5, slice[0].Dummy)
	assert.Equal(t, "subtitle6", slice[0].Title)
}

func TestDeleteMisconfigured(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("DELETE", "/misconfigured/5/subdummy/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
