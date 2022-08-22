package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/diogomattioli/crud/pkg/data"
	"github.com/gorilla/mux"
)

func Create[T data.CreateValidator](w http.ResponseWriter, r *http.Request) {

	var obj T

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	fmt.Fprintf(w, "%v", string(bytes))
}

func Retrieve[T any](w http.ResponseWriter, r *http.Request) {

	var obj T

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := db.First(&obj, id)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%v", string(bytes))
}

func Update[T data.UpdateValidator[T]](w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var obj T

	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var old T

	res := db.First(&old, id)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	fmt.Fprintf(w, "%v", string(bytes))
}

func Delete[T data.DeleteValidator](w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var obj T

	res := db.First(&obj, id)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !obj.IsValidDelete() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	res = db.Delete(obj, id)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
