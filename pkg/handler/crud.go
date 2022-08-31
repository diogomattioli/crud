package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/diogomattioli/crud/pkg/data"
	"github.com/gorilla/mux"
)

func Create[T data.CreateValidator](w http.ResponseWriter, r *http.Request) {

	vars, err := data.VarsInt(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytesVars, _ := json.Marshal(vars)
	var where T
	json.Unmarshal(bytesVars, &where)

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var obj T

	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(bytesVars, &obj)

	if !obj.IsValidCreate() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	res := db.Create(&obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(bytes))
}

func Retrieve[T any](w http.ResponseWriter, r *http.Request) {

	vars, err := data.VarsInt(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytesVars, _ := json.Marshal(vars)
	var where T
	json.Unmarshal(bytesVars, &where)

	var obj T

	res := db.Where(where).Or("1 != 1").First(&obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(bytes))
}

func Update[T data.UpdateValidator[T]](w http.ResponseWriter, r *http.Request) {

	vars, err := data.VarsInt(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytesVars, _ := json.Marshal(vars)
	var where T
	json.Unmarshal(bytesVars, &where)

	var old T

	res := db.Where(where).Or("1 != 1").First(&old)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var obj T = old

	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(bytesVars, &obj)

	if !obj.IsValidUpdate(old) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	res = db.Save(&obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(bytes))
}

func Delete[T data.DeleteValidator](w http.ResponseWriter, r *http.Request) {

	vars, err := data.VarsInt(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var obj T

	bytesVars, _ := json.Marshal(vars)
	var where T
	json.Unmarshal(bytesVars, &where)

	res := db.Where(where).Or("1 != 1").First(&obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !obj.IsValidDelete() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	res = db.Delete(obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
