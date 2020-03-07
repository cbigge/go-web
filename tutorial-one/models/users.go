package models

import (
	"errors"

	"github.com/cbigge/go-web/tutorial-one/hash"
	"github.com/cbigge/go-web/tutorial-one/rand"

	"github.com/jinzhu/gorm"
	// Not directly called but required for postgres to work
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

const (
	hmacSecretKey = "super-secret-hmac-key"
	userPwPepper  = "super-secret-pepper"
)

var (
	// ErrNotFound is returned when a resource cannot be
	// found in the db
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is passed
	// to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invlalid password
	// is entered on authentication
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

// NewUserService creates a db connection
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		db: db,
	}, nil
}

// AutoMigrate will attempt to automatically migrate
// user table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// ByID queries the db to find a User by their
// id and returns the User object
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail queries the db to find a User by their
// email and returns the User object
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember quieres the db to find a User by
// their remember token and returns the User object
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create will create user in db and backfill data
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)

	return us.db.Create(user).Error
}

// Update will update queried user with all data
// provided in the user object parameter
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Authenticate checks user credentials with the provided
// email address and password
//
// If email is invalid,
//   return nil, ErrNotFound
// If password is invalid,
//   return nil, ErrInvalidPassword
// If both email and password are valid,
//   return user, nil
// Otherwise if another error is encountered,
//   return nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}

}

// Close closes the UserService db connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops and rebuilds the user table
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
