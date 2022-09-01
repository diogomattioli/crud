package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEmpty(t *testing.T) {

	setupDb(0)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListNoParams(t *testing.T) {

	setupDb(250)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

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

	setupDb(250)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?offset=50&limit=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 50, len(slice))
	assert.Equal(t, 51, slice[0].ID)
	assert.Equal(t, 100, slice[49].ID)
}

func TestListLimit25(t *testing.T) {

	setupDb(250)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?limit=25", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 25, len(slice))
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 25, slice[24].ID)
}

func TestListFields(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?fields=ID", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(slice))
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, "", slice[0].Title)
	assert.Equal(t, false, slice[0].Valid)

	req, err = http.NewRequest("GET", "/dummy/?fields=Title,Valid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	slice = []Dummy{}
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(slice))
	assert.Equal(t, 0, slice[0].ID)
	assert.Equal(t, "title1", slice[0].Title)
	assert.Equal(t, true, slice[0].Valid)
}

func TestList1Page(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 10, len(slice))
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 10, slice[9].ID)
}

func TestListLimitBadRequest(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?limit=1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?limit=-1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?limit=abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListOffsetBadRequest(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?offset=-1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?offset=abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListSearch(t *testing.T) {

	setupDb(25)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?search=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "3", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

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

	setupDb(25)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?search=5&search=3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "6", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

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

func TestListSort(t *testing.T) {

	setupDb(5)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?sort=Title", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(slice))
	assert.Equal(t, 5, slice[0].ID)
	assert.Equal(t, 1, slice[4].ID)
}

func TestListSubNoParams(t *testing.T) {

	setupDb(250)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/23/subdummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "2", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []SubDummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(slice))
	assert.Equal(t, 45, slice[0].ID)
	assert.Equal(t, 23, slice[0].Dummy)
	assert.Equal(t, 46, slice[1].ID)
	assert.Equal(t, 23, slice[1].Dummy)
}
