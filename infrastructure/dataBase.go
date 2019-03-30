package infrastructure

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	// Read Environment DB Variables
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)

	timeout := time.After(60 * time.Second)

	// Waits until DB is up (Timeout 60 Seconds)
	for {
		select {
		case <-timeout:
			panic(err)
			return
		default:
			db, err = gorm.Open("postgres", dbUri)
			if err == nil {
				// Connected With DB
				return
			}
		}
		time.Sleep(time.Second)
	}

}

// Returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
