package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	DisplayName    string             `json:"displayName,omitempty"`
	Username       string             `json:"username,omitempty"`
	StandardClaims jwt.StandardClaims `json:"-"`
}
