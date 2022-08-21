package handler

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dummy struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}

func setupDb(quantity int) {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Dummy{})

	for i := 1; i <= quantity; i++ {
		db.Create(&Dummy{ID: i, Title: fmt.Sprintf("title%v", quantity-i+1)})
	}

	SetDatabase(db)
}

func destroyDb() {
	db, err := db.DB()
	if err == nil {
		db.Close()
	}
}
