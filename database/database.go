package database

import (
	"github.com/jinzhu/gorm"
	// sqlite
	_ "github.com/mattn/go-sqlite3"
)

// Todo is struct about todo information
type Todo struct {
	gorm.Model
	Todo    string `json:"todo"`
	Process string `json:"process"`
	UserID  string `json:"userid"`
	TodoID  string `json:"todoid"`
}

// User is user information
type User struct {
	gorm.Model
	UserName string `json:"username"`
	UserID   string `json:"userid"`
	Password string `json:"password"`
	Email    string `json:"Email"`
}

var db *gorm.DB

func init() {
	conn, err := gorm.Open("sqlite3", "./data.db")
	if err != nil {
		return
	}
	db = conn
	db.AutoMigrate(&Todo{})
	db.AutoMigrate(&User{})
}

// GetDB return db
func GetDB() *gorm.DB {
	return db
}
