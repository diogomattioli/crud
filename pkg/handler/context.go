package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diogomattioli/crud/pkg/data"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDatabase(_db *gorm.DB) {
	db = _db
}

func getObject[T any](vars []byte) (T, error) {

	var obj T

	var where T

	err := json.Unmarshal(vars, &where)
	if err != nil {
		return obj, errors.New("unmarshal failed")
	}

	res := db.Where(where).Or("1 != 1").First(&obj)
	if res.RowsAffected == 0 {
		return obj, errors.New("object not found")
	}

	return obj, nil
}

func varsToJson(r *http.Request) ([]byte, error) {

	vars, err := data.VarsInt(mux.Vars(r))
	if err != nil {
		return nil, errors.New("invalid vars")
	}

	bytes, err := json.Marshal(vars)
	if err != nil {
		return nil, errors.New("marshal failed")
	}

	return bytes, nil
}
