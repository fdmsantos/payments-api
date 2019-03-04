package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {

	e := godotenv.Load("/home/fsantos/Documents/Projects/Personal/test/.env")
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Exec("DROP TABLE payments;")
	db.Exec("DROP TABLE attributes;")
	db.Exec("DROP TABLE beneficiary_parties;")
	db.Exec("DROP TABLE charges;")
	db.Exec("DROP TABLE charges_informations;")
	db.Exec("DROP TABLE debtor_parties;")
	db.Exec("DROP TABLE fxes;")
	db.Exec("DROP TABLE sponsor_parties;")

	db.AutoMigrate(
		&Account{},
		&Payment{},
		&Attributes{},
		&BeneficiaryParty{},
		&DebtorParty{},
		&SponsorParty{},
		&ChargesInformation{},
		&Charge{},
		&FX{},
	)
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
