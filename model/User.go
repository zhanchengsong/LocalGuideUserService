package model

import (
	"github.com/jinzhu/gorm"
)

type DatabaseError struct {
	ErrorMsg string
	Cause    string
	Code     string
}

func (err DatabaseError) Error() string {
	return err.ErrorMsg
}

type User struct {
	gorm.Model   `json:"-"`
	DisplayName  string `json:"displayName,omitempty"`
	UserId       string `json:"userId,omitempty" gorm:"type:varchar(100);primaryKey"`
	Username     string `json:"username,omitempty" gorm:"type:varchar(100);uniqueIndex"`
	Email        string `json:"email,omitempty" gorm:"type:varchar(100);unique"`
	ImageId      string `json:"imageId,omitempty" gorm:"type:varchar(100)"`
	Password     string `json:"password,omitempty"`
	JWTToken     string `json:"jwtToken,omitempty" gorm:"-"`
	RefreshToken string `json:"refreshToken,omitempty" gorm:"-"`
}
