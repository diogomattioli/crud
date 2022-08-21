package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEmpty(t *testing.T) {

	db := mockDb(0)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListNoParams(t *testing.T) {

	db := mockDb(250)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 50)
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 50, slice[49].ID)
}

func TestList2ndPage(t *testing.T) {

	db := mockDb(250)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?page=2", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "2", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 50)
	assert.Equal(t, 51, slice[0].ID)
	assert.Equal(t, 100, slice[49].ID)
}

func TestList25RecordsPage(t *testing.T) {

	db := mockDb(250)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?records=25", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "25", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 25)
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 25, slice[24].ID)
}

func TestList1Page(t *testing.T) {

	db := mockDb(10)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 10)
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 10, slice[9].ID)
}

func TestListRecordsBadRequest(t *testing.T) {

	db := mockDb(10)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?records=1000", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?records=-1000", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?records=abc", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListPageBadRequest(t *testing.T) {

	db := mockDb(10)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?page=1000", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?page=-1000", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?page=abc", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListIds(t *testing.T) {

	db := mockDb(25)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?id=13&id=19&id=21", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "3", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy

	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 3)
	assert.Equal(t, 13, slice[0].ID)
	assert.Equal(t, 19, slice[1].ID)
	assert.Equal(t, 21, slice[2].ID)
}

func TestListSearch(t *testing.T) {

	db := mockDb(25)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?search=5", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "3", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy

	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 3)
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 11, slice[1].ID)
	assert.Equal(t, 21, slice[2].ID)
}

func TestListSearchMultiple(t *testing.T) {

	db := mockDb(25)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?search=5&search=3", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "6", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy

	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 6)
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 3, slice[1].ID)
	assert.Equal(t, 11, slice[2].ID)
	assert.Equal(t, 13, slice[3].ID)
	assert.Equal(t, 21, slice[4].ID)
	assert.Equal(t, 23, slice[5].ID)
}

func TestListOrder(t *testing.T) {

	db := mockDb(5)
	defer destroyDb(db)

	SetDatabase(db)

	req, err := http.NewRequest("GET", "/dummy/?order=Title", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(List[Dummy])
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Page"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Pages"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "50", rec.Header().Get("X-Paging-RecordsPerPage"))

	var slice []Dummy

	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 5)
	assert.Equal(t, 5, slice[0].ID)
	assert.Equal(t, 1, slice[4].ID)
}
