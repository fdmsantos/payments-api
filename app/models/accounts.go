package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"payments/infrastructure"
	"payments/utils"
	"strings"
	"time"
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

// CreateToken creates a token after a success login
func (a *Account) CreateToken() {
	tk := &Token{
		UserId: a.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		}}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	a.Token, _ = token.SignedString([]byte(os.Getenv("token_password")))
}

// IsEmailValid check if email is valid
func (a *Account) IsEmailValid() error {
	// Check if Email contains @ character
	if !strings.Contains(a.Email, "@") {
		return errors.New(utils.ERROR_EMAIL_REQUIRED)
	}

	// Email must be unique
	tempAccount, err := GetAccountByEmail(a.Email)
	if err != nil {
		return err
	}

	if tempAccount.Email != "" {
		return errors.New(utils.ERROR_EMAIL_ALREADY_EXISTS)
	}

	return nil
}

// IsPasswordValid check if password is valid
func (a *Account) IsPasswordValid() bool {
	// Check if Password has 6 or more characters
	return len(a.Password) >= 6
}

// CreateHashedPassword check if password is valid
func (a *Account) CreateHashedPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err == nil {
		a.Password = string(hashedPassword)
	}
	return err
}

// CheckPassword check if the password is correct to the user
func (a *Account) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return errors.New(utils.ERROR_INVALID_LOGIN)
	}
	return nil
}

// GetAccountByEmail Get a account model through an email
func GetAccountByEmail(email string) (Account, error) {
	account := Account{}
	err := infrastructure.GetDB().Table("accounts").Where("email = ?", email).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return account, errors.New(utils.ERROR_SERVER)
	}
	return account, nil
}
