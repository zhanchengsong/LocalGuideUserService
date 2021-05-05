package transferObject

// Types for request and response for user register

type UserRegisterBody struct {
	DisplayName string `json:"displayName,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	ImageId     string `json:"imageId,omitempty"`
}

type UserUpdateBody struct {
	DisplayName string `json:"displayName,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	ImageId     string `json:"imageId,omitempty"`
}

type UserLoginBody struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserResponseBody struct {
	DisplayName  string `json:"displayName,omitempty"`
	Username     string `json:"username,omitempty"`
	UserId       string `json:"userId,omitempty"`
	Email        string `json:"email,omitempty"`
	JWTToken     string `json:"jwtToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type User struct {
	DisplayName string `json:"displayName,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
}

type UsersResponseBody struct {
	Users []User `json:"users"`
}
