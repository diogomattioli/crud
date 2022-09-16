package handler

import (
	"net/http"

	"github.com/diogomattioli/crud/pkg/data"
)

func sub[T any, S any](w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request)) {
	
	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = getObject[S](vars)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f(w, r)
}

func sub2[T any, S2 any, S any](w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request)) {

	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = getObject[S](vars)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = getObject[S2](vars)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f(w, r)
}

func CreateSub[T data.CreateValidator, S any](w http.ResponseWriter, r *http.Request) {
	sub[T, S](w, r, Create[T])
}

func CreateSub2[T data.CreateValidator, S2 any, S any](w http.ResponseWriter, r *http.Request) {
	sub2[T, S2, S](w, r, Create[T])
}

func RetrieveSub[T any, S any](w http.ResponseWriter, r *http.Request) {
	sub[T, S](w, r, Retrieve[T])
}

func RetrieveSub2[T any, S2 any, S any](w http.ResponseWriter, r *http.Request) {
	sub2[T, S2, S](w, r, Retrieve[T])
}

func UpdateSub[T data.UpdateValidator[T], S any](w http.ResponseWriter, r *http.Request) {
	sub[T, S](w, r, Update[T])
}

func UpdateSub2[T data.UpdateValidator[T], S2 any, S any](w http.ResponseWriter, r *http.Request) {
	sub2[T, S2, S](w, r, Update[T])
}

func DeleteSub[T data.DeleteValidator, S any](w http.ResponseWriter, r *http.Request) {
	sub[T, S](w, r, Delete[T])
}

func DeleteSub2[T data.DeleteValidator, S2 any, S any](w http.ResponseWriter, r *http.Request) {
	sub2[T, S2, S](w, r, Delete[T])
}

func ListSub[T any, S any](w http.ResponseWriter, r *http.Request) {
	sub[T, S](w, r, List[T])
}

func ListSub2[T any, S2 any, S any](w http.ResponseWriter, r *http.Request) {
	sub2[T, S2, S](w, r, List[T])
}