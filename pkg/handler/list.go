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

func createSearchQuery(db *gorm.DB, obj any, query []string) *gorm.DB {

	if len(query) == 0 {
		return db
	}

	var fields []string

	ty := reflect.TypeOf(obj).Elem()
	for i := 0; i < ty.NumField(); i++ {
		if ty.Field(i).Type.Name() == "string" || ty.Field(i).Type.Name() == "NullString" {
			fields = append(fields, ty.Field(i).Name)
		}
	}

	for i := 0; i < len(query); i++ {
		if str := query[i]; data.Valid(str) {
			for _, field := range fields {
				db = db.Or(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", data.ToSnakeCase(field)), "%"+str+"%")
			}
		}
	}

	return db
}

func createSortQuery(db *gorm.DB, obj any, query string) (*gorm.DB, error) {

	if query == "" {
		return db, nil
	}

	if data.Valid(query) {
		if field, exists := reflect.TypeOf(obj).Elem().FieldByName(query); exists {
			return db.Order(fmt.Sprintf("%s ASC", field.Name)), nil
		}
		return db, errors.New("inexistent sort field")
	}

	return db, nil
}

func selectReturnedFields(db *gorm.DB, obj any, query string) (*gorm.DB, error) {

	if query == "" {
		return db, nil
	}

	var fields []string

	queries := strings.Split(strings.ReplaceAll(query, " ", ""), ",")

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

	innerDb, err = selectReturnedFields(innerDb, &obj, URLQuery.Get("fields"))
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

	res := innerDb.Offset(offset).Limit(limit).Where(where).Find(&slice)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(slice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("X-Paging-Total", fmt.Sprint(total))
	w.Header().Add("X-Paging-MaxLimit", fmt.Sprint(maxLimit))

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%v", string(bytes))
}
