package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/diogomattioli/crud/pkg/data"
	"gorm.io/gorm"
)

const (
	defaultRecordsPerPage = 50
	maxRecordsPerPage     = 250
)

func search(db *gorm.DB, obj any, values url.Values) *gorm.DB {

	var fields []string

	ty := reflect.TypeOf(obj).Elem()
	for i := 0; i < ty.NumField(); i++ {
		if ty.Field(i).Type.Name() == "string" || ty.Field(i).Type.Name() == "NullString" {
			fields = append(fields, ty.Field(i).Name)
		}
	}

	for i := 0; i < len(values["search"]); i++ {
		if str := values["search"][i]; data.Valid(str) {
			for _, field := range fields {
				db = db.Or(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", data.ToSnakeCase(field)), "%"+str+"%")
			}
		}
	}

	return db
}

func order(db *gorm.DB, obj any, values url.Values) *gorm.DB {

	if order := values.Get("order"); data.Valid(order) {
		if field, exists := reflect.TypeOf(obj).Elem().FieldByName(order); exists {
			return db.Order(fmt.Sprintf("%s ASC", field.Name))
		}
	}

	return db
}

func List[T any](w http.ResponseWriter, r *http.Request) {

	var slice []T
	var obj T

	innerDb := db

	// Filters
	if ids := r.URL.Query()["id"]; len(ids) > 0 {
		innerDb = innerDb.Or(ids)
	}
	innerDb = search(innerDb, &obj, r.URL.Query())
	innerDb = order(innerDb, &obj, r.URL.Query())
	// Filters

	var total int64
	innerDb.Model(obj).Count(&total)
	if total == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var err error

	recordsPerPage := defaultRecordsPerPage
	if data.Valid(r.URL.Query().Get("records")) {
		recordsPerPage, err = strconv.Atoi(r.URL.Query().Get("records"))
		if err != nil || recordsPerPage < 1 || recordsPerPage > maxRecordsPerPage {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	pages := int(total / int64(recordsPerPage))
	if total%int64(recordsPerPage) > 0 {
		pages++
	}

	page := 1
	if data.Valid(r.URL.Query().Get("page")) {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 || page > pages {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	res := innerDb.Offset((page - 1) * recordsPerPage).Limit(recordsPerPage).Find(&slice)
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(slice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("X-Paging-Page", fmt.Sprint(page))
	w.Header().Add("X-Paging-Pages", fmt.Sprint(pages))
	w.Header().Add("X-Paging-Total", fmt.Sprint(total))
	w.Header().Add("X-Paging-RecordsPerPage", fmt.Sprint(recordsPerPage))
	w.Header().Add("X-Paging-MaxRecordsPerPage", fmt.Sprint(maxRecordsPerPage))

	fmt.Fprintf(w, "%v", string(bytes))
}
