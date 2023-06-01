package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/diogomattioli/crud/pkg/data"
)

func Create[T data.CreateValidator](w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var obj T

	// unmarshall the object from body
	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// overwrite id with provided in the vars/url
	err = json.Unmarshal(vars, &obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), Session{}, Session{Token: r.Header.Get("X-Access-Token")})

	err = obj.ValidateCreate(ctx)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "%v", err)
		return
	}

	res := db.Create(&obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%+v%+v", r.URL.RequestURI(), obj.GetID()))
	w.Header().Set("X-Item-ID", fmt.Sprintf("%+v", obj.GetID()))

	w.WriteHeader(http.StatusCreated)
}

func Retrieve[T any](w http.ResponseWriter, r *http.Request) {

	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	obj, err := getObject[T](vars)
	if err != nil {
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

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	old, err := getObject[T](vars)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var obj T = old

	// unmarshall the object from body
	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// overwrite id with provided in the vars/url
	err = json.Unmarshal(vars, &obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), Session{}, Session{Token: r.Header.Get("X-Access-Token")})

	err = obj.ValidateUpdate(ctx, old)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "%v", err)
		return
	}

	res := db.Save(&obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func Delete[T data.DeleteValidator](w http.ResponseWriter, r *http.Request) {

	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	obj, err := getObject[T](vars)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.WithValue(r.Context(), Session{}, Session{Token: r.Header.Get("X-Access-Token")})

	err = obj.ValidateDelete(ctx)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "%v", err)
		return
	}

	res := db.Delete(obj)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
