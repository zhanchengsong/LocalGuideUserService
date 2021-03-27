package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	DisplayName string `json:"displayName,omitempty"`
	Username    string `json:"username" gorm:"type:varchar(100);unique_index"`
	Email       string `json:"email" gorm:"type:varchar(100);unique_index"`
	Password    string `json:"password"`
	IconUrl     string `json:"iconUrl,omitempty"`
	JWTToken    string `json:"jwtToken,omitempty"`
}
