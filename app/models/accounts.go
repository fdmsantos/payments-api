package models

import (
	//"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}
