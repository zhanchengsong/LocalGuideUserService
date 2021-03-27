package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	DisplayName string `json:"displayName,omitempty"`
	UserId      string `json:"-" gorm:"type:varchar(100);unique_index"`
	Username    string `json:"username",omitempty gorm:"type:varchar(100);unique_index"`
	Email       string `json:"email",omitempty gorm:"type:varchar(100);unique_index"`
	Password    string `json:"-",omitempty`
	IconUrl     string `json:"iconUrl,omitempty"`
	JWTToken    string `json:"jwtToken,omitempty"`
}
