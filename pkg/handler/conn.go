package handler

import "gorm.io/gorm"

var db *gorm.DB

func SetDatabase(_db *gorm.DB) {
	db = _db
}
