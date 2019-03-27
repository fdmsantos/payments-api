package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

var db *gorm.DB

func init() {

	// Read Environment DB Variables
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	for {
		// Waits from DB is UP
		dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
		conn, err := gorm.Open("postgres", dbUri)
		if err != nil {
			fmt.Print(err)
		} else {
			db = conn
			break
		}
	}

}

// Returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
