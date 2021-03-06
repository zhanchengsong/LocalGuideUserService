// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/countByDisplayName": {
            "get": {
                "description": "Return total count of users with the given display name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Count user by display name"
                ],
                "summary": "Get count of users in the database by displayName",
                "parameters": [
                    {
                        "type": "string",
                        "description": "exact display name string to search",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.UserCountResponseBody"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    }
                }
            }
        },
        "/countByUsername": {
            "get": {
                "description": "Return total count of users with the given username",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Count user by username"
                ],
                "summary": "Get count of users in the database by username",
                "parameters": [
                    {
                        "type": "string",
                        "description": "exact username string to search",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.UserCountResponseBody"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Takes in username and password to assign token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Log in a user"
                ],
                "summary": "Login a user and obtain jwtToken/refreshToken",
                "parameters": [
                    {
                        "description": "Login user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transferObject.UserLoginBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transferObject.UserResponseBody"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    }
                }
            }
        },
        "/token": {
            "post": {
                "description": "Use referesh token to obtain new jwt token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Refresh JWT Token"
                ],
                "summary": "Referesh JWT Token using the refresh token",
                "parameters": [
                    {
                        "description": "RefreshToken",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenResponseBody"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "Create a user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Create a user"
                ],
                "summary": "Create a user",
                "parameters": [
                    {
                        "description": "Create user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transferObject.UserRegisterBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/transferObject.UserResponseBody"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Search for users with username that is partial matching the string",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Find user by partial matching"
                ],
                "summary": "Get a list of users with username that contains the string provided",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Term for partial matching username",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transferObject.UsersResponseBody"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.handlerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.TokenRequestBody": {
            "type": "object",
            "properties": {
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "handlers.TokenResponseBody": {
            "type": "object",
            "properties": {
                "jwtToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "handlers.UserCountResponseBody": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                }
            }
        },
        "handlers.handlerError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        },
        "transferObject.User": {
            "type": "object",
            "properties": {
                "displayName": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transferObject.UserLoginBody": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transferObject.UserRegisterBody": {
            "type": "object",
            "properties": {
                "displayName": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transferObject.UserResponseBody": {
            "type": "object",
            "properties": {
                "displayName": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "jwtToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transferObject.UsersResponseBody": {
            "type": "object",
            "properties": {
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/transferObject.User"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "auth.zhancheng.dev",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "User API",
	Description: "Service for user registration, login and jwtToken refresh",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
