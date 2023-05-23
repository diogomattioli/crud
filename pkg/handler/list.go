package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/diogomattioli/crud/pkg/data"
	"gorm.io/gorm"
)

const (
	maxLimit     = 250
	defaultLimit = 50
)

func createSearchQuery[T any](db *gorm.DB, obj T, queries []string) *gorm.DB {

	ty := reflect.TypeOf(obj).Elem()

	for _, query := range queries {
		for j := 0; j < ty.NumField(); j++ {

			typeName := ty.Field(j).Type.Name()
			fieldName := ty.Field(j).Name

			if data.Valid(query) && (typeName == "string" || typeName == "NullString") {
				db = db.Or(fmt.Sprintf("%s LIKE LOWER(?)", data.ToSnakeCase(fieldName)), "%"+query+"%")
			} else if value, err := strconv.Atoi(query); err == nil && (strings.HasPrefix(typeName, "int") || strings.HasPrefix(typeName, "NullInt")) {
				db = db.Or(fmt.Sprintf("%s = ?", data.ToSnakeCase(fieldName)), value)
			} else if value, err := strconv.ParseFloat(query, 64); err == nil && (strings.HasPrefix(typeName, "float") || strings.HasPrefix(typeName, "NullFloat")) {
				db = db.Or(fmt.Sprintf("%s = ?", data.ToSnakeCase(fieldName)), value)
			}
		}
	}

	return db
}

func createSortQuery[T any](db *gorm.DB, obj T, query string) (*gorm.DB, error) {

	if query == "" {
		return db, nil
	}

	if data.Valid(query) {
		if field, exists := reflect.TypeOf(obj).Elem().FieldByName(query); exists {
			return db.Order(fmt.Sprintf("%s ASC", data.ToSnakeCase(field.Name))), nil
		}
		return db, errors.New("inexistent sort field")
	}

	return db, nil
}

func selectReturnedFields[T any](db *gorm.DB, obj T, queries []string) (*gorm.DB, error) {
	var fields []string

	ty := reflect.TypeOf(obj).Elem()
	for _, query := range queries {
		for i := 0; i < ty.NumField(); i++ {
			if ty.Field(i).Name == query {
				fields = append(fields, ty.Field(i).Name)
				break
			}
		}
	}

	if len(fields) < len(queries) {
		return db, errors.New("inexistent field")
	}

	if len(fields) > 0 {
		return db.Select(fields), nil
	}

	return db, nil
}

func List[T any](w http.ResponseWriter, r *http.Request) {

	vars, err := varsToJson(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var where T

	err = json.Unmarshal(vars, &where)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("X-Paging-MaxLimit", fmt.Sprint(maxLimit))
	w.Header().Add("X-Paging-DefaultLimit", fmt.Sprint(defaultLimit))

	var slice []T
	var obj T

	innerDb := db

	URLQuery := r.URL.Query()

	offset := 0
	if URLQuery.Get("offset") != "" {
		offset, err = strconv.Atoi(URLQuery.Get("offset"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	limit := defaultLimit
	if URLQuery.Get("limit") != "" {
		limit, err = strconv.Atoi(URLQuery.Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if offset < 0 || limit <= 0 || limit > maxLimit {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	innerDb, err = selectReturnedFields(innerDb, &obj, URLQuery["field"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Filters
	innerDb = createSearchQuery(innerDb, &obj, URLQuery["search"])
	innerDb, err = createSortQuery(innerDb, &obj, URLQuery.Get("sort"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Filters

	var total int64
	innerDb.Model(obj).Where(where).Count(&total)
	if total == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("X-Paging-Total", fmt.Sprint(total))

	innerDb.Offset(offset).Limit(limit).Where(where).Find(&slice)
	if len(slice) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("X-Paging-Size", fmt.Sprint(len(slice)))

	bytes, err := json.Marshal(slice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%v", string(bytes))
}
