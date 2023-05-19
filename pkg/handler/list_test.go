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

	setupDb(5)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Size"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(slice), 5)
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 5, slice[4].ID)
}

func TestList2ndPage(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?offset=5&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Size"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(slice))
	assert.Equal(t, 6, slice[0].ID)
	assert.Equal(t, 10, slice[4].ID)
}

func TestListLimit5(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Size"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(slice))
	assert.Equal(t, 1, slice[0].ID)
	assert.Equal(t, 5, slice[4].ID)
}

func TestListOffset(t *testing.T) {

	setupDb(10)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?offset=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Size"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(slice))
	assert.Equal(t, 6, slice[0].ID)
	assert.Equal(t, 10, slice[4].ID)
}

func TestListLimitOffset(t *testing.T) {

	setupDb(15)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?offset=5&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "15", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Size"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []Dummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(slice))
	assert.Equal(t, 6, slice[0].ID)
	assert.Equal(t, 10, slice[4].ID)
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
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Size"))
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
	assert.Equal(t, "1", rec.Header().Get("X-Paging-Size"))
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

func TestListFieldsWrongField(t *testing.T) {

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?fields=IDWrong", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("GET", "/dummy/?fields=TitleWrong,Valid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
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
	assert.Equal(t, "10", rec.Header().Get("X-Paging-Size"))
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

	setupDb(0)
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

	setupDb(0)
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
	assert.Equal(t, "3", rec.Header().Get("X-Paging-Size"))
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
	assert.Equal(t, "6", rec.Header().Get("X-Paging-Size"))
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
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Size"))
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

func TestListSortID(t *testing.T) {

	setupDb(5)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy_default/?sort=DummyDefaultID", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "5", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []DummyDefault
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(slice))
	assert.Equal(t, 1, slice[0].DummyDefaultID)
	assert.Equal(t, 5, slice[4].DummyDefaultID)
}

func TestListSortWrongField(t *testing.T) {

	setupDb(5)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/?sort=TitleWrong", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListSubNoParams(t *testing.T) {

	setupDb(5)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/3/subdummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "2", rec.Header().Get("X-Paging-Total"))
	assert.Equal(t, "2", rec.Header().Get("X-Paging-Size"))
	assert.Equal(t, "250", rec.Header().Get("X-Paging-MaxLimit"))

	var slice []SubDummy
	err = json.NewDecoder(rec.Body).Decode(&slice)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(slice))
	assert.Equal(t, 5, slice[0].ID)
	assert.Equal(t, 3, slice[0].Dummy)
	assert.Equal(t, 6, slice[1].ID)
	assert.Equal(t, 3, slice[1].Dummy)
}

func TestListSubNotFound(t *testing.T) {

	setupDb(2)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/dummy/23/subdummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListMisconfigured(t *testing.T) {

	setupDb(2)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/misconfigured/23/subdummy/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTP(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
