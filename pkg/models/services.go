package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//NewUserService : Instantiates a new service to deal with users.
func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db: db}
	uv := &userValidiator{UserDB: ug}
	usc := &userService{UserDB: uv}
	return usc
}

//UserService : An interface that decides the higher level functionality
//of dealing with users.
type UserService interface {
	Authenticate(user *User) error
	UserDB
}

type userService struct {
	UserDB
}

func (us *userService) Authenticate(user *User) error {
	//check if the user is actually valid
	founduser, err := us.ByEmail(user.Email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(founduser.Password), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

//UserValidator : We will create structs that fit this interface.
//The functions that this struct will recieve will clean up
//and do basic validation.
type UserValidator interface {
	UserDB
}

type userValidiator struct {
	UserDB
}

func (uv *userValidiator) Insert(user *User) error {
	//check is user email is empty
	if user.Email == "" {
		return errors.New("Error: Email Empty")
	}

	newUser, _ := uv.UserDB.ByEmail(user.Email)

	if newUser != nil {
		return errors.New("Error: Account Exists")
	}

	if user.Password == "" {
		return errors.New("Error: Password empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Error: Error Validating Password")
	}
	user.Password = string(hash)
	return uv.UserDB.Insert(user)
}

//UserDB : An interface that decides the contracts for the methods provided
//to interact with the database.
type UserDB interface {
	TableExists() bool
	AutoMigrate() error
	ByID(id int) (*User, error)
	ByEmail(email string) (*User, error)
	Insert(*User) error
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id int) (*User, error) {
	var user User
	err := ug.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	err := ug.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Insert(user *User) error {
	err := ug.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) AutoMigrate() error {
	if ug.TableExists() {
		return nil
	}

	err := ug.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) TableExists() bool {
	return ug.db.HasTable(&User{})
}
